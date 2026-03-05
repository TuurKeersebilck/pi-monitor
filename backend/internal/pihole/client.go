package pihole

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	baseURL  string
	password string
	sid      string
	sidMu    sync.Mutex
	http     *http.Client
}

type Stats struct {
	TotalQueries   int     `json:"total_queries"`
	BlockedQueries int     `json:"blocked_queries"`
	BlockedPercent float64 `json:"blocked_percent"`
	DomainsBlocked int     `json:"domains_blocked"`
}

type authResponse struct {
	Session struct {
		SID string `json:"sid"`
	} `json:"session"`
}

type summaryResponse struct {
	Queries struct {
		Total   int     `json:"total"`
		Blocked int     `json:"blocked"`
		PercentBlocked float64 `json:"percent_blocked"`
	} `json:"queries"`
	Gravity struct {
		DomainsBeingBlocked int `json:"domains_being_blocked"`
	} `json:"gravity"`
}

func NewClient(baseURL, password string) *Client {
	return &Client{
		baseURL:  baseURL,
		password: password,
		http:     &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) authenticate() error {
	body, _ := json.Marshal(map[string]string{"password": c.password})
	resp, err := c.http.Post(c.baseURL+"/api/auth", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth failed with status %d", resp.StatusCode)
	}

	var authResp authResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("auth decode failed: %w", err)
	}

	c.sidMu.Lock()
	c.sid = authResp.Session.SID
	c.sidMu.Unlock()

	return nil
}

func (c *Client) GetStats() (*Stats, error) {
	c.sidMu.Lock()
	sid := c.sid
	c.sidMu.Unlock()

	// Authenticate if we don't have a session
	if sid == "" {
		if err := c.authenticate(); err != nil {
			return nil, err
		}
		c.sidMu.Lock()
		sid = c.sid
		c.sidMu.Unlock()
	}

	req, err := http.NewRequest("GET", c.baseURL+"/api/stats/summary", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("sid", sid)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Session expired — re-auth once and retry
	if resp.StatusCode == http.StatusUnauthorized {
		c.sidMu.Lock()
		c.sid = ""
		c.sidMu.Unlock()

		if err := c.authenticate(); err != nil {
			return nil, err
		}

		c.sidMu.Lock()
		sid = c.sid
		c.sidMu.Unlock()

		req, _ = http.NewRequest("GET", c.baseURL+"/api/stats/summary", nil)
		req.Header.Set("sid", sid)
		resp, err = c.http.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}

	var summary summaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		return nil, fmt.Errorf("decode summary: %w", err)
	}

	return &Stats{
		TotalQueries:   summary.Queries.Total,
		BlockedQueries: summary.Queries.Blocked,
		BlockedPercent: summary.Queries.PercentBlocked,
		DomainsBlocked: summary.Gravity.DomainsBeingBlocked,
	}, nil
}
