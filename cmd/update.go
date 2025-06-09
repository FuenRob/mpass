package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mpass/internal/storage"
	"mpass/internal/ui"
	"time"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a password entry",
	Long:  "Update a password entry with new details such as username, URL, or password",
	RunE:  runUpdate,
}

func runUpdate(_ *cobra.Command, _ []string) error {
	masterPassword, err := ui.PromptPassword("Enter master password:")
	if err != nil {
		return fmt.Errorf("failed to get master password: %w", err)
	}

	vaultManager := storage.NewVault()

	entries, err := vaultManager.GetAllEntries(masterPassword)
	if err != nil {
		return fmt.Errorf("failed to load entries: %w", err)
	}
	if len(entries) == 0 {
		return fmt.Errorf("No entries found in vault")
	}

	selectedEntry, err := ui.SelectEntry(entries)
	if err != nil {
		return fmt.Errorf("failed to select entry: %w", err)
	}

	fmt.Println("Leave any field blank to keep it unchanged.")
	newUsername, _ := ui.PromptInput("New Username (current: " + selectedEntry.Username + "):")
	newURL, _ := ui.PromptInput("New URL (current: " + selectedEntry.URL + "):")
	newPassword, _ := ui.PromptPassword("New Password (leave blank so as not to change it):")

	updated := false
	for i := range entries {
		if entries[i].Username == selectedEntry.Username && entries[i].URL == selectedEntry.URL && entries[i].Password == selectedEntry.Password {
			if newUsername != "" {
				entries[i].Username = newUsername
				updated = true
			}
			if newURL != "" {
				entries[i].URL = newURL
				updated = true
			}
			if newPassword != "" {
				entries[i].Password = newPassword
				updated = true
			}
			if updated {
				entries[i].UpdatedAt = time.Now()
			}
			break
		}
	}
	if !updated {
		fmt.Println("No changes were made.")
		return nil
	}

	err = vaultManager.UpdateEntries(entries, masterPassword)
	if err != nil {
		return fmt.Errorf("failed to save updated entry: %w", err)
	}
	fmt.Printf("✅ Password updated for %s copied to clipboard!\n",
		selectedEntry.URL)
	return nil
}
