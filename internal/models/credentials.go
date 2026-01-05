// internal/models/credential.go
package models

import "time"

type Credential struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Username          string    `json:"username"`
	EncryptedPassword string    `json:"encrypted_password"`
	URL               string    `json:"url"`
	Notes             string    `json:"notes"`
	Category          string    `json:"category"`
	Created           time.Time `json:"created"`
	Updated           time.Time `json:"updated"`
}

type VaultConfig struct {
	ID           string `json:"id"`
	Salt         string `json:"salt"`
	PasswordHash string `json:"password_hash"`
	Created      string `json:"created"`
}
