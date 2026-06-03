package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/tuurk/dashboard/internal/api"
	"github.com/tuurk/dashboard/internal/config"
	"github.com/tuurk/dashboard/internal/docker"
	"github.com/tuurk/dashboard/internal/immich"
	"github.com/tuurk/dashboard/internal/pihole"
	"github.com/tuurk/dashboard/internal/store"
	"github.com/tuurk/dashboard/internal/system"
	"github.com/tuurk/dashboard/internal/ws"
)

const writeInterval = 10 * time.Second

func main() {
	piholeURL := getenv("PIHOLE_URL", "http://pihole:80")
	piholePassword := getenv("PIHOLE_APP_PASSWORD", "")
	immichURL := getenv("IMMICH_URL", "")
	immichAPIKey := getenv("IMMICH_API_KEY", "")
	staticDir := getenv("STATIC_DIR", "/app/frontend/dist")
	listenAddr := getenv("LISTEN_ADDR", ":8080")
	dataDir := getenv("DATA_DIR", "/app/data")

	cfgStore, err := config.NewStore(dataDir)
	if err != nil {
		log.Fatalf("failed to init config store: %v", err)
	}

	uploadsDir := filepath.Join(dataDir, "uploads")
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Fatalf("failed to create uploads dir: %v", err)
	}

	// Stats persistence — non-fatal if SQLite fails (e.g. read-only volume)
	var statsStore *store.DB
	statsStore, err = store.Open(filepath.Join(dataDir, "stats.db"))
	if err != nil {
		log.Printf("warning: stats store unavailable: %v — persistence disabled", err)
	} else {
		log.Println("stats store opened")
		cfg, _ := cfgStore.Load()
		statsStore.Prune(retentionDays(cfg))
		startWriter(statsStore, cfgStore)
	}

	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Fatalf("failed to create docker client: %v", err)
	}
	defer dockerClient.Close()

	piholeClient := pihole.NewClient(piholeURL, piholePassword)

	var immichClient *immich.Client
	if immichURL != "" && immichAPIKey != "" {
		immichClient = immich.NewClient(immichURL, immichAPIKey)
		log.Printf("immich integration enabled: %s", immichURL)
	} else {
		log.Println("immich integration disabled (IMMICH_URL or IMMICH_API_KEY not set)")
	}

	hub := ws.NewHub(dockerClient, piholeClient, immichClient)

	handler := api.NewHandler(hub, cfgStore, statsStore, uploadsDir)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, staticDir)

	log.Printf("dashboard listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// startWriter starts the always-on background goroutines that write stats to disk
// and prune old records, independent of whether any browser is connected.
func startWriter(statsStore *store.DB, cfgStore *config.Store) {
	// Separate Docker client so its internal network-delta state doesn't
	// interfere with the hub's client (they poll at different intervals).
	writerDocker, err := docker.NewClient()
	if err != nil {
		log.Printf("writer: docker client failed, container stats won't be persisted: %v", err)
	}

	go func() {
		ticker := time.NewTicker(writeInterval)
		defer ticker.Stop()
		for range ticker.C {
			ts := time.Now().Unix()

			if stats, err := system.GetStats(); err == nil {
				statsStore.WriteSystemStat(
					ts,
					stats.CPU.UsagePercent,
					stats.RAM.UsagePercent,
					int64(stats.RAM.UsedMB),
					stats.Disk.UsagePercent,
					stats.Temp.CPUTempC,
					stats.Net.RxBytesPerSec,
					stats.Net.TxBytesPerSec,
				)
			}

			if writerDocker != nil {
				ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
				ctrs, err := writerDocker.ListContainers(ctx)
				cancel()
				if err == nil {
					var inputs []store.ContainerStatInput
					for _, c := range ctrs {
						if c.Running {
							inputs = append(inputs, store.ContainerStatInput{
								Name:        c.Name,
								CPUPercent:  c.CPUPercent,
								MemPercent:  c.MemPercent,
								MemUsedMB:   c.MemUsedMB,
								NetRxBytesS: c.NetRxBytesS,
								NetTxBytesS: c.NetTxBytesS,
							})
						}
					}
					if len(inputs) > 0 {
						statsStore.WriteContainerStats(ts, inputs)
					}
				}
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			cfg, _ := cfgStore.Load()
			statsStore.Prune(retentionDays(cfg))
		}
	}()
}

func retentionDays(cfg *config.Config) int {
	if cfg == nil || cfg.RetentionDays <= 0 {
		return 7
	}
	return cfg.RetentionDays
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
