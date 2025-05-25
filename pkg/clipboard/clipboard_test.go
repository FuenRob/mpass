package clipboard

import (
	"os/exec"
	"runtime"
	"testing"
)

func TestWriteText(t *testing.T) {
	testText := "test-clipboard-content"

	// Check if clipboard tools are available
	if !isClipboardAvailable() {
		t.Skip("Clipboard tools not available, skipping test")
	}

	err := WriteText(testText)
	if err != nil {
		t.Fatalf("Failed to write to clipboard: %v", err)
	}

	// Note: Reading from clipboard is platform-specific and complex
	// This test only verifies that WriteText doesn't return an error
	// Manual verification would be needed to ensure content is actually copied
}

func TestWriteEmptyText(t *testing.T) {
	if !isClipboardAvailable() {
		t.Skip("Clipboard tools not available, skipping test")
	}

	err := WriteText("")
	if err != nil {
		t.Fatalf("Failed to write empty text to clipboard: %v", err)
	}
}

func TestWriteMultilineText(t *testing.T) {
	if !isClipboardAvailable() {
		t.Skip("Clipboard tools not available, skipping test")
	}

	multilineText := "line1\nline2\nline3"
	err := WriteText(multilineText)
	if err != nil {
		t.Fatalf("Failed to write multiline text to clipboard: %v", err)
	}
}

func TestWriteSpecialCharacters(t *testing.T) {
	if !isClipboardAvailable() {
		t.Skip("Clipboard tools not available, skipping test")
	}

	specialText := "Special chars: !@#$%^&*()_+{}|:<>?[]\\;'\",./"
	err := WriteText(specialText)
	if err != nil {
		t.Fatalf("Failed to write special characters to clipboard: %v", err)
	}
}

// isClipboardAvailable checks if the required clipboard tools are available
func isClipboardAvailable() bool {
	switch runtime.GOOS {
	case "darwin":
		_, err := exec.LookPath("pbcopy")
		return err == nil
	case "linux":
		// Check for either xclip or xsel
		_, err1 := exec.LookPath("xclip")
		_, err2 := exec.LookPath("xsel")
		return err1 == nil || err2 == nil
	case "windows":
		_, err := exec.LookPath("clip")
		return err == nil
	default:
		return false
	}
}

func TestTryCommand(t *testing.T) {
	// Test con un comando que siempre debe existir
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/C", "echo", "test"}
	default:
		cmd = "echo"
		args = []string{"test"}
	}

	err := tryCommand(cmd, args, "input")
	if err != nil {
		t.Fatalf("tryCommand falló con un comando básico: %v", err)
	}
}

func TestTryCommandNonexistent(t *testing.T) {
	err := tryCommand("nonexistent-command-12345", []string{}, "input")
	if err == nil {
		t.Fatal("tryCommand should fail with nonexistent command")
	}
}
