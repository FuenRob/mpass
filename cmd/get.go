package cmd

import (
	"fmt"
	"mpass/internal/models"
	"mpass/internal/storage"
	"mpass/internal/ui"
	"mpass/pkg/clipboard"

	"github.com/spf13/cobra"
)

var (
	getCmd = &cobra.Command{
		Use:   "get",
		Short: "Get a password entry",
		Long:  "Search and retrieve a password entry by username or URL",
		RunE:  runGet,
	}
	searchUser string
	searchURL  string
)

// init initializes the flags for the getCmd command.
// It sets up the command-line options for searching by username or URL.
func init() {
	getCmd.Flags().StringVarP(&searchUser, "user", "u", "", "Search by username")
	getCmd.Flags().StringVarP(&searchURL, "url", "l", "", "Search by URL")
}

// runGet executes the logic for the "get" command.
// It prompts the user for the master password, searches for password entries
// by username or URL, allows selection if multiple entries are found, and
// copies the selected password to the clipboard.
func runGet(_ *cobra.Command, _ []string) error {
	if searchUser == "" && searchURL == "" {
		return fmt.Errorf("please provide either --user or --url flag")
	}

	// Get master password
	masterPassword, err := ui.PromptPassword("Enter master password:")
	if err != nil {
		return fmt.Errorf("failed to get master password: %w", err)
	}

	// Load vault
	vault := storage.NewVault()
	entries, err := vault.SearchEntries(searchUser, searchURL, masterPassword)
	if err != nil {
		return fmt.Errorf("failed to search entries: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("❌ No matching entries found")
		return nil
	}

	// If multiple entries, let user select
	var selectedEntry *models.PasswordEntry
	if len(entries) == 1 {
		selectedEntry = &entries[0]
	} else {
		selected, err := ui.SelectEntry(entries)
		if err != nil {
			return fmt.Errorf("failed to select entry: %w", err)
		}
		selectedEntry = selected
	}

	// Copy password to clipboard
	if err := clipboard.WriteText(selectedEntry.Password); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	fmt.Printf("✅ Password for %s@%s copied to clipboard!\n",
		selectedEntry.Username, selectedEntry.URL)
	return nil
}
