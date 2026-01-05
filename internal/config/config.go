// internal/config/config.go
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	PocketBaseURL   string `json:"pocketbase_url"`
	AdminEmail      string `json:"admin_email"`
	SessionTimeout  int    `json:"session_timeout_minutes"`
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".passmanager", "config.json")
}

func Load() (*Config, error) {
	configPath := GetConfigPath()
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Save() error {
	configPath := GetConfigPath()
	
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}
