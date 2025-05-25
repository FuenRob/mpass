package storage

import (
	"mpass/internal/models"
	"os"
	"path/filepath"
	"testing"
)

func createTestVault(t *testing.T) (*VaultManager, string) {
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test_vault.enc")
	return &VaultManager{vaultPath: vaultPath}, tempDir
}

func TestNewVault(t *testing.T) {
	vault := NewVault()
	if vault == nil {
		t.Fatal("NewVault() returned nil")
	}

	homeDir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homeDir, ".mpass", "vault.enc")
	if vault.vaultPath != expectedPath {
		t.Fatalf("Expected vault path %s, got %s", expectedPath, vault.vaultPath)
	}
}

func TestAddAndGetEntry(t *testing.T) {
	vault, _ := createTestVault(t)
	masterPassword := "test-master-password"

	entry := models.PasswordEntry{
		Username: "testuser",
		URL:      "https://example.com",
		Password: "secret123",
	}

	// Add entry
	err := vault.AddEntry(entry, masterPassword)
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	// Get all entries
	entries, err := vault.GetAllEntries(masterPassword)
	if err != nil {
		t.Fatalf("Failed to get entries: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}

	retrievedEntry := entries[0]
	if retrievedEntry.Username != entry.Username {
		t.Fatalf("Expected username %s, got %s", entry.Username, retrievedEntry.Username)
	}
	if retrievedEntry.URL != entry.URL {
		t.Fatalf("Expected URL %s, got %s", entry.URL, retrievedEntry.URL)
	}
	if retrievedEntry.Password != entry.Password {
		t.Fatalf("Expected password %s, got %s", entry.Password, retrievedEntry.Password)
	}

	// Verify timestamps are set
	if retrievedEntry.CreatedAt.IsZero() {
		t.Fatal("CreatedAt should be set")
	}
	if retrievedEntry.UpdatedAt.IsZero() {
		t.Fatal("UpdatedAt should be set")
	}
}

func TestAddMultipleEntries(t *testing.T) {
	vault, _ := createTestVault(t)
	masterPassword := "test-master-password"

	entries := []models.PasswordEntry{
		{Username: "user1", URL: "https://site1.com", Password: "pass1"},
		{Username: "user2", URL: "https://site2.com", Password: "pass2"},
		{Username: "user3", URL: "https://site3.com", Password: "pass3"},
	}

	// Add all entries
	for _, entry := range entries {
		err := vault.AddEntry(entry, masterPassword)
		if err != nil {
			t.Fatalf("Failed to add entry: %v", err)
		}
	}

	// Retrieve all entries
	retrievedEntries, err := vault.GetAllEntries(masterPassword)
	if err != nil {
		t.Fatalf("Failed to get entries: %v", err)
	}

	if len(retrievedEntries) != len(entries) {
		t.Fatalf("Expected %d entries, got %d", len(entries), len(retrievedEntries))
	}

	// Verify all entries are present
	for i, expected := range entries {
		found := false
		for _, retrieved := range retrievedEntries {
			if retrieved.Username == expected.Username &&
				retrieved.URL == expected.URL &&
				retrieved.Password == expected.Password {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Entry %d not found in retrieved entries", i)
		}
	}
}

func TestSearchEntries(t *testing.T) {
	vault, _ := createTestVault(t)
	masterPassword := "test-master-password"

	entries := []models.PasswordEntry{
		{Username: "john@gmail.com", URL: "https://github.com", Password: "pass1"},
		{Username: "john@yahoo.com", URL: "https://gitlab.com", Password: "pass2"},
		{Username: "mary@gmail.com", URL: "https://github.com", Password: "pass3"},
		{Username: "admin", URL: "https://company.com", Password: "pass4"},
	}

	// Add all entries
	for _, entry := range entries {
		err := vault.AddEntry(entry, masterPassword)
		if err != nil {
			t.Fatalf("Failed to add entry: %v", err)
		}
	}

	// Test search by username
	results, err := vault.SearchEntries("john", "", masterPassword)
	if err != nil {
		t.Fatalf("Failed to search by username: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("Expected 2 results for 'john', got %d", len(results))
	}

	// Test search by URL
	results, err = vault.SearchEntries("", "github", masterPassword)
	if err != nil {
		t.Fatalf("Failed to search by URL: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("Expected 2 results for 'github', got %d", len(results))
	}

	// Test search by both username and URL
	results, err = vault.SearchEntries("john", "github", masterPassword)
	if err != nil {
		t.Fatalf("Failed to search by username and URL: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 result for 'john' and 'github', got %d", len(results))
	}
	if results[0].Username != "john@gmail.com" || results[0].URL != "https://github.com" {
		t.Fatal("Wrong entry returned for combined search")
	}

	// Test search with no matches
	results, err = vault.SearchEntries("nonexistent", "", masterPassword)
	if err != nil {
		t.Fatalf("Failed to search for nonexistent username: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("Expected 0 results for nonexistent user, got %d", len(results))
	}
}

func TestWrongMasterPassword(t *testing.T) {
	vault, _ := createTestVault(t)
	correctPassword := "correct-password"
	wrongPassword := "wrong-password"

	entry := models.PasswordEntry{
		Username: "testuser",
		URL:      "https://example.com",
		Password: "secret123",
	}

	// Add entry with correct password
	err := vault.AddEntry(entry, correctPassword)
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	// Try to retrieve with wrong password
	_, err = vault.GetAllEntries(wrongPassword)
	if err == nil {
		t.Fatal("Should fail with wrong master password")
	}

	// Try to search with wrong password
	_, err = vault.SearchEntries("testuser", "", wrongPassword)
	if err == nil {
		t.Fatal("Should fail with wrong master password")
	}
}

func TestEmptyVault(t *testing.T) {
	vault, _ := createTestVault(t)
	masterPassword := "test-password"

	// Get entries from empty vault
	entries, err := vault.GetAllEntries(masterPassword)
	if err != nil {
		t.Fatalf("Failed to get entries from empty vault: %v", err)
	}

	if len(entries) != 0 {
		t.Fatalf("Expected 0 entries from empty vault, got %d", len(entries))
	}

	// Search in empty vault
	results, err := vault.SearchEntries("anything", "", masterPassword)
	if err != nil {
		t.Fatalf("Failed to search in empty vault: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("Expected 0 results from empty vault, got %d", len(results))
	}
}

func TestVaultPersistence(t *testing.T) {
	vault, _ := createTestVault(t)
	masterPassword := "test-password"

	entry := models.PasswordEntry{
		Username: "testuser",
		URL:      "https://example.com",
		Password: "secret123",
	}

	// Add entry
	err := vault.AddEntry(entry, masterPassword)
	if err != nil {
		t.Fatalf("Failed to add entry: %v", err)
	}

	// Create new vault instance with same path
	vault2 := &VaultManager{vaultPath: vault.vaultPath}

	// Retrieve entries with new instance
	entries, err := vault2.GetAllEntries(masterPassword)
	if err != nil {
		t.Fatalf("Failed to get entries with new vault instance: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}

	if entries[0].Username != entry.Username {
		t.Fatal("Entry not persisted correctly")
	}
}
