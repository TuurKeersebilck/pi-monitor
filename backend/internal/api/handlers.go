package api

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tuurk/dashboard/internal/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Safe since we're LAN only behind Caddy
		return true
	},
}

type Handler struct {
	hub *ws.Hub
}

func NewHandler(hub *ws.Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, staticDir string) {
	// WebSocket endpoint
	mux.HandleFunc("/ws", h.handleWS)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Serve Svelte frontend static files
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
}

func (h *Handler) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws: upgrade error: %v", err)
		return
	}

	log.Printf("ws: client connected from %s", r.RemoteAddr)
	h.hub.AddClient(conn)

	// Keep connection open, remove client when it closes
	go func() {
		defer h.hub.RemoveClient(conn)
		for {
			// Read messages — we don't use them but need this to detect disconnects
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Printf("ws: client disconnected: %v", err)
				return
			}
		}
	}()
}
