package config

import (
	"testing"

	"passmanager/internal/models"
)

func TestNewDefault(t *testing.T) {
	cfg := NewDefault()
	if cfg.Initialized {
		t.Fatal("expected Initialized to be false by default")
	}
	if cfg.Settings == nil {
		t.Fatal("expected default Settings to be non-nil")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	if Exists() {
		t.Fatal("expected Exists() to be false before saving")
	}

	cfg := &Config{
		PocketBaseURL: "http://127.0.0.1:8090",
		AdminEmail:    "admin@example.com",
		Initialized:   true,
	}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if !Exists() {
		t.Fatal("expected Exists() to be true after saving")
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.PocketBaseURL != cfg.PocketBaseURL {
		t.Fatalf("got PocketBaseURL %q, want %q", loaded.PocketBaseURL, cfg.PocketBaseURL)
	}
	if loaded.AdminEmail != cfg.AdminEmail {
		t.Fatalf("got AdminEmail %q, want %q", loaded.AdminEmail, cfg.AdminEmail)
	}
	if !loaded.Initialized {
		t.Fatal("expected Loaded.Initialized to be true")
	}
	if loaded.Settings == nil {
		t.Fatal("expected Loaded.Settings to be populated with defaults")
	}
	if loaded.Settings.SessionTimeout != models.DefaultSettings().SessionTimeout {
		t.Fatalf("got SessionTimeout %d, want %d", loaded.Settings.SessionTimeout, models.DefaultSettings().SessionTimeout)
	}
}
