package models

import "time"

// PasswordEntry represents a single password entry
type PasswordEntry struct {
	Username  string    `json:"username"`
	URL       string    `json:"url"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Vault represents the encrypted storage container
type Vault struct {
	Entries []PasswordEntry `json:"entries"`
	Salt    []byte          `json:"salt"`
}
