// internal/config/config.go
package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"passmanager/internal/models"
)

type Config struct {
	PocketBaseURL string              `json:"pocketbase_url"`
	AdminEmail    string              `json:"admin_email"`
	Settings      *models.AppSettings `json:"settings"`
	Initialized   bool                `json:"initialized"`
}

func GetConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".passmanager")
}

func GetConfigPath() string {
	return filepath.Join(GetConfigDir(), "config.json")
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

	// Set defaults if settings are nil
	if config.Settings == nil {
		config.Settings = models.DefaultSettings()
	}

	return &config, nil
}

func (c *Config) Save() error {
	configDir := GetConfigDir()
	configPath := GetConfigPath()

	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	if c.Settings == nil {
		c.Settings = models.DefaultSettings()
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

func Exists() bool {
	_, err := os.Stat(GetConfigPath())
	return err == nil
}

func NewDefault() *Config {
	return &Config{
		Settings:    models.DefaultSettings(),
		Initialized: false,
	}
}