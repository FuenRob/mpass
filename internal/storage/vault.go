package storage

import (
	"encoding/json"
	"fmt"
	"mpass/internal/crypto"
	"mpass/internal/models"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	vaultDir  = ".mpass"
	vaultFile = "vault.enc"
)

type VaultManager struct {
	vaultPath string
}

// NewVault creates a new VaultManager instance with the default vault path
// located in the user's home directory.
func NewVault() *VaultManager {
	homeDir, _ := os.UserHomeDir()
	vaultPath := filepath.Join(homeDir, vaultDir, vaultFile)
	return &VaultManager{vaultPath: vaultPath}
}

// ensureVaultDir creates the directory for the vault file if it does not exist.
// It returns an error if the directory cannot be created.
func (v *VaultManager) ensureVaultDir() error {
	dir := filepath.Dir(v.vaultPath)
	return os.MkdirAll(dir, 0700)
}

// loadVault loads the encrypted vault from disk, decrypts it using the provided master password,
// and returns the Vault model. If the vault file does not exist, it creates a new vault with a random salt.
// Returns an error if reading, decrypting, or parsing the vault fails.
func (v *VaultManager) loadVault(masterPassword string) (*models.Vault, error) {
	if _, err := os.Stat(v.vaultPath); os.IsNotExist(err) {
		// Create new vault with random salt
		salt, err := crypto.GenerateSalt()
		if err != nil {
			return nil, fmt.Errorf("failed to generate salt: %w", err)
		}
		return &models.Vault{
			Entries: []models.PasswordEntry{},
			Salt:    salt,
		}, nil
	}

	// Load existing vault
	encryptedData, err := os.ReadFile(v.vaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}

	// First, we need to extract the salt from the beginning of the file
	if len(encryptedData) < 32 {
		return nil, fmt.Errorf("invalid vault file format")
	}

	salt := encryptedData[:32]
	ciphertext := encryptedData[32:]

	// Derive key and decrypt
	key := crypto.DeriveKey(masterPassword, salt)
	decryptedData, err := crypto.Decrypt(ciphertext, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault (wrong password?): %w", err)
	}

	var vault models.Vault
	if err := json.Unmarshal(decryptedData, &vault); err != nil {
		return nil, fmt.Errorf("failed to parse vault data: %w", err)
	}

	vault.Salt = salt
	return &vault, nil
}

// saveVault serializes the given Vault struct, encrypts it using the provided master password,
// and writes the encrypted data (with prepended salt) to disk. It ensures the vault directory exists
// and returns an error if any step fails.
func (v *VaultManager) saveVault(vault *models.Vault, masterPassword string) error {
	if err := v.ensureVaultDir(); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	// Serialize vault data
	data, err := json.Marshal(vault)
	if err != nil {
		return fmt.Errorf("failed to serialize vault: %w", err)
	}

	// Derive key and encrypt
	key := crypto.DeriveKey(masterPassword, vault.Salt)
	encryptedData, err := crypto.Encrypt(data, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault: %w", err)
	}

	// Prepend salt to encrypted data
	finalData := append(vault.Salt, encryptedData...)

	// Write to file with secure permissions
	if err := os.WriteFile(v.vaultPath, finalData, 0600); err != nil {
		return fmt.Errorf("failed to write vault file: %w", err)
	}

	return nil
}

// AddEntry adds a new password entry to the vault, setting the creation and update timestamps,
// and saves the updated vault encrypted with the provided master password.
// Returns an error if loading or saving the vault fails.
func (v *VaultManager) AddEntry(entry models.PasswordEntry, masterPassword string) error {
	vault, err := v.loadVault(masterPassword)
	if err != nil {
		return err
	}

	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	vault.Entries = append(vault.Entries, entry)

	return v.saveVault(vault, masterPassword)
}

// GetAllEntries loads the vault using the provided master password and returns all password entries.
// Returns a slice of PasswordEntry and an error if loading the vault fails.
func (v *VaultManager) GetAllEntries(masterPassword string) ([]models.PasswordEntry, error) {
	vault, err := v.loadVault(masterPassword)
	if err != nil {
		return nil, err
	}
	return vault.Entries, nil
}

// SearchEntries searches for password entries in the vault that match the given username and/or URL.
// It loads the vault using the provided master password and performs a case-insensitive substring match
// on the username and URL fields. If either parameter is an empty string, it is ignored in the search.
// Returns a slice of matching PasswordEntry structs and an error if loading the vault fails.
func (v *VaultManager) SearchEntries(username, url, masterPassword string) ([]models.PasswordEntry, error) {
	vault, err := v.loadVault(masterPassword)
	if err != nil {
		return nil, err
	}

	var matches []models.PasswordEntry
	for _, entry := range vault.Entries {
		usernameMatch := username == "" || strings.Contains(strings.ToLower(entry.Username), strings.ToLower(username))
		urlMatch := url == "" || strings.Contains(strings.ToLower(entry.URL), strings.ToLower(url))

		if usernameMatch && urlMatch {
			matches = append(matches, entry)
		}
	}

	return matches, nil
}
