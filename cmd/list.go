package cmd

import (
	"fmt"
	"mpass/internal/storage"
	"mpass/internal/ui"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all password entries",
	Long:  "Display all stored password entries (passwords are hidden)",
	RunE:  runList,
}

// runList executes the "list" command, prompting the user for the master password,
// loading all password entries from the vault, and displaying them (without showing passwords).
// Returns an error if the master password is not provided or if entries cannot be loaded.
func runList(_ *cobra.Command, _ []string) error {
	// Get master password
	masterPassword, err := ui.PromptPassword("Enter master password:")
	if err != nil {
		return fmt.Errorf("failed to get master password: %w", err)
	}

	// Load vault
	vault := storage.NewVault()
	entries, err := vault.GetAllEntries(masterPassword)
	if err != nil {
		return fmt.Errorf("failed to load entries: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("📭 No password entries found")
		return nil
	}

	fmt.Printf("📚 Found %d password entries:\n\n", len(entries))
	for i, entry := range entries {
		fmt.Printf("%d. %s@%s\n", i+1, entry.Username, entry.URL)
	}

	return nil
}
