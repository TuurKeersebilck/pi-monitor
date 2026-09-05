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
	fastPollInterval      = 1 * time.Second  // system stats
	containerPollInterval = 3 * time.Second  // docker container stats (expensive: one Docker API call per container)
	slowPollInterval      = 30 * time.Second // pi-hole + immich (overview data)
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

	containerPollMu   sync.Mutex
	lastContainerPoll time.Time
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
	go h.broadcastContainersThrottled()
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

	go h.pollEvery(ctx, fastPollInterval, h.broadcastFast)
	go h.pollEvery(ctx, containerPollInterval, h.broadcastContainersThrottled)
	go h.pollEvery(ctx, slowPollInterval, h.broadcastSlow)
}

// pollEvery calls fn on every tick of interval until ctx is cancelled.
func (h *Hub) pollEvery(ctx context.Context, interval time.Duration, fn func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fn()
		}
	}
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
	stats, err := system.GetStats()
	if err != nil {
		log.Printf("ws: system stats error: %v", err)
		return
	}
	h.send(Payload{System: stats})
}

// broadcastContainersThrottled runs broadcastContainers unless the last run
// was less than containerPollInterval ago. Container stats are expensive
// (one Docker API call per container), and both the poll ticker and every
// new client connection call this -- without throttling, a flaky connection
// or several tabs reconnecting fires far more of these scans than the
// ticker alone ever would, which is exactly the "hammering" this interval
// was introduced to stop.
func (h *Hub) broadcastContainersThrottled() {
	h.containerPollMu.Lock()
	if time.Since(h.lastContainerPoll) < containerPollInterval {
		h.containerPollMu.Unlock()
		return
	}
	h.lastContainerPoll = time.Now()
	h.containerPollMu.Unlock()

	h.broadcastContainers()
}

func (h *Hub) broadcastContainers() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	containers, err := h.docker.ListContainers(ctx)
	if err != nil {
		log.Printf("ws: docker error: %v", err)
		return
	}
	h.send(Payload{Containers: containers})
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
