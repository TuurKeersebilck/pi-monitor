package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tuurk/dashboard/internal/api"
	"github.com/tuurk/dashboard/internal/docker"
	"github.com/tuurk/dashboard/internal/immich"
	"github.com/tuurk/dashboard/internal/pihole"
	"github.com/tuurk/dashboard/internal/ws"
)

func main() {
	piholeURL := getenv("PIHOLE_URL", "http://pihole:80")
	piholePassword := getenv("PIHOLE_APP_PASSWORD", "")
	immichURL := getenv("IMMICH_URL", "")
	immichAPIKey := getenv("IMMICH_API_KEY", "")
	staticDir := getenv("STATIC_DIR", "/app/frontend/dist")
	listenAddr := getenv("LISTEN_ADDR", ":8080")

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

	handler := api.NewHandler(hub)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, staticDir)

	log.Printf("dashboard listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
