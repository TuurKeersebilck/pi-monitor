package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tuurk/dashboard/internal/docker"
	"github.com/tuurk/dashboard/internal/immich"
	"github.com/tuurk/dashboard/internal/pihole"
	"github.com/tuurk/dashboard/internal/system"
)

const (
	fastPollInterval = 1 * time.Second  // system stats + containers
	slowPollInterval = 30 * time.Second // pi-hole + immich (overview data)
)

type Payload struct {
	System     *system.Stats         `json:"system,omitempty"`
	Containers []docker.ContainerInfo `json:"containers,omitempty"`
	Pihole     *pihole.Stats         `json:"pihole,omitempty"`
	Immich     *immich.Stats         `json:"immich,omitempty"`
	Error      string                `json:"error,omitempty"`
}

type Hub struct {
	clients    map[*websocket.Conn]struct{}
	mu         sync.Mutex
	docker     *docker.Client
	pihole     *pihole.Client
	immich     *immich.Client
	cancelPoll context.CancelFunc
}

func NewHub(dockerClient *docker.Client, piholeClient *pihole.Client, immichClient *immich.Client) *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]struct{}),
		docker:  dockerClient,
		pihole:  piholeClient,
		immich:  immichClient,
	}
}

func (h *Hub) AddClient(conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[conn] = struct{}{}
	wasEmpty := len(h.clients) == 1
	h.mu.Unlock()

	if wasEmpty {
		h.startPolling()
	}

	go h.broadcastFast()
	go h.broadcastSlow()
}

func (h *Hub) RemoveClient(conn *websocket.Conn) {
	h.mu.Lock()
	delete(h.clients, conn)
	isEmpty := len(h.clients) == 0
	h.mu.Unlock()

	conn.Close()

	if isEmpty && h.cancelPoll != nil {
		log.Println("ws: no clients connected, stopping poll")
		h.cancelPoll()
		h.cancelPoll = nil
	}
}

func (h *Hub) startPolling() {
	log.Println("ws: client connected, starting poll")
	ctx, cancel := context.WithCancel(context.Background())
	h.cancelPoll = cancel

	go func() {
		ticker := time.NewTicker(fastPollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				h.broadcastFast()
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(slowPollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				h.broadcastSlow()
			}
		}
	}()
}

func (h *Hub) send(payload Payload) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ws: marshal error: %v", err)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("ws: write error, removing client: %v", err)
			delete(h.clients, conn)
			conn.Close()
		}
	}
}

func (h *Hub) broadcastFast() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		stats, err := system.GetStats()
		if err != nil {
			log.Printf("ws: system stats error: %v", err)
		} else {
			h.send(Payload{System: stats})
		}
	}()

	go func() {
		defer wg.Done()
		containers, err := h.docker.ListContainers(ctx)
		if err != nil {
			log.Printf("ws: docker error: %v", err)
		} else {
			h.send(Payload{Containers: containers})
		}
	}()

	wg.Wait()
}

func (h *Hub) broadcastSlow() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		stats, err := h.pihole.GetStats()
		if err != nil {
			log.Printf("ws: pihole error: %v", err)
		} else {
			h.send(Payload{Pihole: stats})
		}
	}()

	if h.immich != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stats, err := h.immich.GetStats()
			if err != nil {
				log.Printf("ws: immich error: %v", err)
			} else {
				h.send(Payload{Immich: stats})
			}
		}()
	}

	wg.Wait()
}
