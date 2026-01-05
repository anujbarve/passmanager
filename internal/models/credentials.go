// internal/models/credential.go
package models

type Credential struct {
	ID                string `json:"id,omitempty"`
	Title             string `json:"title"`
	Username          string `json:"username,omitempty"`
	EncryptedPassword string `json:"encrypted_password"`
	URL               string `json:"url,omitempty"`
	Notes             string `json:"notes,omitempty"`
	Category          string `json:"category,omitempty"`
	Created           string `json:"created,omitempty"`
	Updated           string `json:"updated,omitempty"`
}

type VaultConfig struct {
	ID           string `json:"id,omitempty"`
	Salt         string `json:"salt"`
	PasswordHash string `json:"password_hash"`
	Created      string `json:"created,omitempty"`
	Updated      string `json:"updated,omitempty"`
}

type AppSettings struct {
	SessionTimeout   int    `json:"session_timeout_minutes"`
	ClipboardTimeout int    `json:"clipboard_timeout_seconds"`
	DefaultCategory  string `json:"default_category"`
	PasswordLength   int    `json:"password_length"`
	IncludeSymbols   bool   `json:"include_symbols"`
}

func DefaultSettings() *AppSettings {
	return &AppSettings{
		SessionTimeout:   5,
		ClipboardTimeout: 30,
		DefaultCategory:  "general",
		PasswordLength:   20,
		IncludeSymbols:   true,
	}
}