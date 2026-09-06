package immich

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tuurk/dashboard/internal/roundutil"
)

type Stats struct {
	Photos   int64   `json:"photos"`
	Videos   int64   `json:"videos"`
	UsageGB  float64 `json:"usage_gb"`
	Healthy  bool    `json:"healthy"`
}

type Client struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		http:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) GetStats() (*Stats, error) {
	healthy := c.ping()

	req, err := http.NewRequest("GET", c.baseURL+"/api/server/statistics", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("immich API returned %d", resp.StatusCode)
	}

	var data struct {
		Photos int64 `json:"photos"`
		Videos int64 `json:"videos"`
		Usage  int64 `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &Stats{
		Photos:  data.Photos,
		Videos:  data.Videos,
		UsageGB: roundutil.Round(float64(data.Usage)/1024/1024/1024, 1),
		Healthy: healthy,
	}, nil
}

func (c *Client) ping() bool {
	req, err := http.NewRequest("GET", c.baseURL+"/api/server/ping", nil)
	if err != nil {
		return false
	}
	req.Header.Set("x-api-key", c.apiKey)
	resp, err := c.http.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
