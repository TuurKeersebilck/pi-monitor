package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Service struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon"`
}

type Config struct {
	Services   []Service `json:"services"`
	Background string    `json:"background"`
}

type Store struct {
	path string
	mu   sync.RWMutex
}

func NewStore(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	return &Store{path: filepath.Join(dataDir, "config.json")}, nil
}

func (s *Store) Load() (*Config, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return &Config{Services: []Service{}}, nil
	}
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Services == nil {
		cfg.Services = []Service{}
	}
	return &cfg, nil
}

func (s *Store) Save(cfg *Config) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}
