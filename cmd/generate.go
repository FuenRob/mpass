package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"mpass/pkg/clipboard"
	"strings"

	"github.com/spf13/cobra"
)

var (
	length      int
	charset     string
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a random password",
		Long:  "Generate a secure random password with specified length and character set",
		RunE:  runGenerate,
	}
)

func init() {
	generateCmd.Flags().IntVarP(&length, "length", "n", 16, "Length of the password to generate")
	generateCmd.Flags().StringVarP(&charset, "charset", "c", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()", "Character set to use for password generation")
}

func runGenerate(_ *cobra.Command, _ []string) error {
	if length <= 0 {
		return fmt.Errorf("Length must be greater than zero")
	}

	if charset == "" {
		return fmt.Errorf("The character set cannot be empty")
	}

	var sb strings.Builder
	var password string

	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return fmt.Errorf("Error generating the password: %w", err)
		}
		sb.WriteByte(charset[idx.Int64()])
	}

	password = sb.String()

	fmt.Println("New password generated:", password)

	if err := clipboard.WriteText(password); err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	fmt.Println("✅ Password copied to clipboard!")

	return nil
}
