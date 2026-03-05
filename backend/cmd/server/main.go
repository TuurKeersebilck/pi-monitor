package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tuurk/dashboard/internal/api"
	"github.com/tuurk/dashboard/internal/docker"
	"github.com/tuurk/dashboard/internal/pihole"
	"github.com/tuurk/dashboard/internal/ws"
)

func main() {
	// Config from environment
	piholeURL := getenv("PIHOLE_URL", "http://pihole:80")
	piholePassword := getenv("PIHOLE_APP_PASSWORD", "")
	staticDir := getenv("STATIC_DIR", "/app/frontend/dist")
	listenAddr := getenv("LISTEN_ADDR", ":8080")

	// Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Fatalf("failed to create docker client: %v", err)
	}
	defer dockerClient.Close()

	// Pi-hole client
	piholeClient := pihole.NewClient(piholeURL, piholePassword)

	// WebSocket hub
	hub := ws.NewHub(dockerClient, piholeClient)

	// HTTP handler
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
