package models

import (
	"testing"
	"time"
)

func TestPasswordEntry(t *testing.T) {
	entry := PasswordEntry{
		Username:  "testuser",
		URL:       "https://example.com",
		Password:  "secret123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if entry.Username != "testuser" {
		t.Fatalf("Expected username 'testuser', got '%s'", entry.Username)
	}

	if entry.URL != "https://example.com" {
		t.Fatalf("Expected URL 'https://example.com', got '%s'", entry.URL)
	}

	if entry.Password != "secret123" {
		t.Fatalf("Expected password 'secret123', got '%s'", entry.Password)
	}

	if entry.CreatedAt.IsZero() {
		t.Fatal("CreatedAt should not be zero")
	}

	if entry.UpdatedAt.IsZero() {
		t.Fatal("UpdatedAt should not be zero")
	}
}

func TestVault(t *testing.T) {
	entries := []PasswordEntry{
		{Username: "user1", URL: "https://site1.com", Password: "pass1"},
		{Username: "user2", URL: "https://site2.com", Password: "pass2"},
	}

	salt := []byte("test-salt-32-bytes")

	vault := Vault{
		Entries: entries,
		Salt:    salt,
	}

	if len(vault.Entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(vault.Entries))
	}

	if string(vault.Salt) != "test-salt-32-bytes" {
		t.Fatalf("Expected salt 'test-salt-32-bytes', got '%s'", string(vault.Salt))
	}

	// Test that entries are accessible
	if vault.Entries[0].Username != "user1" {
		t.Fatalf("Expected first entry username 'user1', got '%s'", vault.Entries[0].Username)
	}

	if vault.Entries[1].Username != "user2" {
		t.Fatalf("Expected second entry username 'user2', got '%s'", vault.Entries[1].Username)
	}
}

func TestPasswordEntryDefaults(t *testing.T) {
	entry := PasswordEntry{}

	if entry.Username != "" {
		t.Fatalf("Expected empty username, got '%s'", entry.Username)
	}

	if entry.URL != "" {
		t.Fatalf("Expected empty URL, got '%s'", entry.URL)
	}

	if entry.Password != "" {
		t.Fatalf("Expected empty password, got '%s'", entry.Password)
	}

	if !entry.CreatedAt.IsZero() {
		t.Fatal("Expected zero CreatedAt")
	}

	if !entry.UpdatedAt.IsZero() {
		t.Fatal("Expected zero UpdatedAt")
	}
}
