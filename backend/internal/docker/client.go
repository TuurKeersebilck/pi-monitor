package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	Status  string `json:"status"`
	Running bool   `json:"running"`
	// UpdateAvailable will be filled in by the registry checker later
	UpdateAvailable bool   `json:"update_available"`
	LatestVersion   string `json:"latest_version,omitempty"`
}

type Client struct {
	cli *client.Client
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) ListContainers(ctx context.Context) ([]ContainerInfo, error) {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var result []ContainerInfo
	for _, ctr := range containers {
		name := "unknown"
		if len(ctr.Names) > 0 {
			// Docker prefixes names with "/"
			name = strings.TrimPrefix(ctr.Names[0], "/")
		}

		result = append(result, ContainerInfo{
			Name:    name,
			Image:   ctr.Image,
			Status:  ctr.Status,
			Running: ctr.State == "running",
		})
	}

	return result, nil
}

func (c *Client) Close() {
	c.cli.Close()
}
