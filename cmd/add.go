package cmd

import (
	"fmt"
	"mpass/internal/models"
	"mpass/internal/storage"
	"mpass/internal/ui"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new password entry",
	Long:  "Add a new password entry with username, URL, and password",
	RunE:  runAdd,
}

func runAdd(_ *cobra.Command, _ []string) error {
	// Get master password
	masterPassword, err := ui.PromptPassword("Enter master password:")
	if err != nil {
		return fmt.Errorf("failed to get master password: %w", err)
	}

	// Get entry details
	username, err := ui.PromptInput("Username:")
	if err != nil {
		return fmt.Errorf("failed to get username: %w", err)
	}

	url, err := ui.PromptInput("URL:")
	if err != nil {
		return fmt.Errorf("failed to get URL: %w", err)
	}

	password, err := ui.PromptPassword("Password:")
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	// Create entry
	entry := models.PasswordEntry{
		Username: username,
		URL:      url,
		Password: password,
	}

	// Save entry
	vault := storage.NewVault()
	if err := vault.AddEntry(entry, masterPassword); err != nil {
		return fmt.Errorf("failed to add entry: %w", err)
	}

	fmt.Println("✅ Password entry added successfully!")
	return nil
}
