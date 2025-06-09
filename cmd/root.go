package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mpass",
	Short: "A secure password manager CLI",
	Long: `Manager Passwords is a secure command-line password manager that stores
your passwords encrypted locally on your machine.`,
}

// Execute runs the root command for the CLI application.
// It handles command execution and prints errors to stderr if any occur.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

// init initializes the root command by adding subcommands to it.
// This function is automatically called when the package is initialized.
func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(generateCmd)
}
