package docker

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	Name            string  `json:"name"`
	Image           string  `json:"image"`
	Status          string  `json:"status"`
	Uptime          string  `json:"uptime"`
	Running         bool    `json:"running"`
	UpdateAvailable bool    `json:"update_available"`
	LatestVersion   string  `json:"latest_version,omitempty"`
	Group           string  `json:"group"`
	URL             string  `json:"url"`
	CPUPercent      float64 `json:"cpu_percent"`
	MemUsedMB       float64 `json:"mem_used_mb"`
	MemLimitMB      float64 `json:"mem_limit_mb"`
	MemPercent      float64 `json:"mem_percent"`
	NetRxBytesS     float64 `json:"net_rx_bytes_s"`
	NetTxBytesS     float64 `json:"net_tx_bytes_s"`
}

// containerStatsJSON holds only the fields we need from Docker's stats API response.
type containerStatsJSON struct {
	CPUStats struct {
		CPUUsage struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemCPUUsage uint64 `json:"system_cpu_usage"`
		OnlineCPUs     int    `json:"online_cpus"`
	} `json:"cpu_stats"`
	PreCPUStats struct {
		CPUUsage struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemCPUUsage uint64 `json:"system_cpu_usage"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Usage uint64            `json:"usage"`
		Limit uint64            `json:"limit"`
		Stats map[string]uint64 `json:"stats"`
	} `json:"memory_stats"`
	Networks map[string]struct {
		RxBytes uint64 `json:"rx_bytes"`
		TxBytes uint64 `json:"tx_bytes"`
	} `json:"networks"`
}

type Client struct {
	cli         *client.Client
	mu          sync.Mutex
	prevNet     map[string][2]uint64 // containerID -> [totalRx, totalTx]
	prevNetTime time.Time
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli, prevNet: make(map[string][2]uint64)}, nil
}

const maxStatsWorkers = 8

func (c *Client) ListContainers(ctx context.Context) ([]ContainerInfo, error) {
	ctrs, err := c.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	c.mu.Lock()
	elapsed := 0.0
	if !c.prevNetTime.IsZero() {
		elapsed = now.Sub(c.prevNetTime).Seconds()
	}
	c.mu.Unlock()

	results := make([]ContainerInfo, len(ctrs))
	sem := make(chan struct{}, maxStatsWorkers)
	var wg sync.WaitGroup

	for i, ctr := range ctrs {
		wg.Add(1)
		i, ctr := i, ctr
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			name := "unknown"
			if len(ctr.Names) > 0 {
				name = strings.TrimPrefix(ctr.Names[0], "/")
			}

			info := ContainerInfo{
				Name:    name,
				Image:   ctr.Image,
				Status:  ctr.Status,
				Uptime:  parseUptime(ctr.Status),
				Running: ctr.State == "running",
				Group:   ctr.Labels["dashboard.group"],
				URL:     ctr.Labels["dashboard.url"],
			}

			if ctr.State == "running" {
				c.fillStats(ctx, &info, ctr.ID, elapsed)
			}

			results[i] = info
		}()
	}

	wg.Wait()

	// Rebuild prevNet with only containers still present so stale IDs don't accumulate.
	c.mu.Lock()
	newPrevNet := make(map[string][2]uint64, len(ctrs))
	for _, ctr := range ctrs {
		if v, ok := c.prevNet[ctr.ID]; ok {
			newPrevNet[ctr.ID] = v
		}
	}
	c.prevNet = newPrevNet
	c.prevNetTime = now
	c.mu.Unlock()

	return results, nil
}

func (c *Client) fillStats(ctx context.Context, info *ContainerInfo, id string, elapsed float64) {
	resp, err := c.cli.ContainerStats(ctx, id, false)
	if err != nil {
		log.Printf("docker: stats unavailable for %s: %v", info.Name, err)
		return
	}
	defer resp.Body.Close()

	var s containerStatsJSON
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return
	}

	// CPU — Docker provides both current and previous snapshot in a single non-streaming call
	cpuDelta := float64(s.CPUStats.CPUUsage.TotalUsage - s.PreCPUStats.CPUUsage.TotalUsage)
	sysDelta := float64(s.CPUStats.SystemCPUUsage - s.PreCPUStats.SystemCPUUsage)
	numCPUs := s.CPUStats.OnlineCPUs
	if numCPUs == 0 {
		numCPUs = 1
	}
	if sysDelta > 0 {
		info.CPUPercent = round(cpuDelta/sysDelta*float64(numCPUs)*100, 1)
	}

	// Memory — subtract inactive file cache (correct for both cgroups v1 and v2)
	memUsed := s.MemoryStats.Usage
	if v, ok := s.MemoryStats.Stats["inactive_file"]; ok && v < memUsed {
		memUsed -= v
	} else if v, ok := s.MemoryStats.Stats["cache"]; ok && v < memUsed {
		memUsed -= v
	}
	info.MemUsedMB = round(float64(memUsed)/1024/1024, 1)
	info.MemLimitMB = round(float64(s.MemoryStats.Limit)/1024/1024, 1)
	if s.MemoryStats.Limit > 0 {
		info.MemPercent = round(float64(memUsed)/float64(s.MemoryStats.Limit)*100, 1)
	}

	// Network — sum across all interfaces, compute bytes/sec from stored previous snapshot
	var totalRx, totalTx uint64
	for _, iface := range s.Networks {
		totalRx += iface.RxBytes
		totalTx += iface.TxBytes
	}

	c.mu.Lock()
	prev, hasPrev := c.prevNet[id]
	c.prevNet[id] = [2]uint64{totalRx, totalTx}
	c.mu.Unlock()

	if hasPrev && elapsed > 0 {
		info.NetRxBytesS = round(float64(totalRx-prev[0])/elapsed, 0)
		info.NetTxBytesS = round(float64(totalTx-prev[1])/elapsed, 0)
	}
}

func (c *Client) Close() {
	c.cli.Close()
}

func parseUptime(status string) string {
	if strings.HasPrefix(status, "Up ") {
		return strings.TrimPrefix(status, "Up ")
	}
	return "stopped"
}

func round(val float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Round(val*pow) / pow
}
