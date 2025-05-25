package ui

import (
	"bufio"
	"fmt"
	"github.com/manifoldco/promptui"
	"golang.org/x/term"
	"mpass/internal/models"
	"os"
	"strings"
	"syscall"
)

// PromptInput prompts the user for input with the given label.
// It returns the trimmed input string or an error if reading fails.
func PromptInput(label string) (string, error) {
	fmt.Print(label + " ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// PromptPassword prompts the user for a password input with the given label.
// The input is hidden (not echoed to the terminal).
// Returns the entered password as a string, or an error if reading fails.
func PromptPassword(label string) (string, error) {
	fmt.Print(label + " ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add newline after hidden input
	if err != nil {
		return "", err
	}
	return string(password), nil
}

// SelectEntry displays a list of PasswordEntry items and allows the user to select one.
// It returns a pointer to the selected PasswordEntry or an error if the selection fails.
func SelectEntry(entries []models.PasswordEntry) (*models.PasswordEntry, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▶ {{ .Username }}@{{ .URL }}",
		Inactive: "  {{ .Username }}@{{ .URL }}",
		Selected: "✅ Selected: {{ .Username }}@{{ .URL }}",
	}

	prompt := promptui.Select{
		Label:     "Multiple entries found. Please select one:",
		Items:     entries,
		Templates: templates,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &entries[index], nil
}
