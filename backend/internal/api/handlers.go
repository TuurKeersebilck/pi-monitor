package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tuurk/dashboard/internal/config"
	"github.com/tuurk/dashboard/internal/ws"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Safe since we're LAN only behind Caddy
		return true
	},
}

type Handler struct {
	hub        *ws.Hub
	cfgStore   *config.Store
	uploadsDir string
}

func NewHandler(hub *ws.Hub, cfgStore *config.Store, uploadsDir string) *Handler {
	return &Handler{hub: hub, cfgStore: cfgStore, uploadsDir: uploadsDir}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, staticDir string) {
	mux.HandleFunc("/ws", h.handleWS)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/config", h.handleConfig)
	mux.HandleFunc("/api/upload", h.handleUpload)

	// Serve uploaded files
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(h.uploadsDir))))

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

	go func() {
		defer h.hub.RemoveClient(conn)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Printf("ws: client disconnected: %v", err)
				return
			}
		}
	}()
}

func (h *Handler) handleConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		cfg, err := h.cfgStore.Load()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(cfg)

	case http.MethodPost:
		var cfg config.Config
		if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.cfgStore.Save(&cfg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(64 << 20); err != nil { // 64 MB limit
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	dst, err := os.Create(filepath.Join(h.uploadsDir, filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": "/uploads/" + filename})
}
