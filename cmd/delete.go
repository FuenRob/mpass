package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mpass/internal/storage"
	"mpass/internal/ui"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a password entry",
	Long:  "Delete a password entry from the vault by selecting it from a list.",
	RunE:  runDelete,
}

func runDelete(cmd *cobra.Command, args []string) error {
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
		return fmt.Errorf("no entries found in vault")
	}

	selectedEntry, err := ui.SelectEntry(entries)
	if err != nil {
		return fmt.Errorf("failed to select entry: %w", err)
	}

	if err := vaultManager.DeleteEntry(selectedEntry, masterPassword); err != nil {
		return fmt.Errorf("failed to delete entry: %w", err)
	}

	fmt.Println("✅ Entry deleted successfully")
	return nil
}
