package clipboard

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// tryCommand executes an external command with the specified arguments and input.
// name: name of the command to execute.
// args: arguments for the command.
// input: text to be passed as standard input to the command.
// Returns an error if the command execution fails.
func tryCommand(name string, args []string, input string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = strings.NewReader(input)
	return cmd.Run()
}

// WriteText copies the provided text to the system clipboard.
// It supports macOS, Linux, and Windows by using the appropriate clipboard utility for each OS.
// Returns an error if the operation fails or the OS is unsupported.
func WriteText(text string) error {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(text)
		return cmd.Run()
	case "linux":
		if err := tryCommand("xclip", []string{"-selection", "clipboard"}, text); err == nil {
			return nil
		}
		return tryCommand("xsel", []string{"--clipboard", "--input"}, text)
	case "windows":
		cmd := exec.Command("clip")
		cmd.Stdin = strings.NewReader(text)
		return cmd.Run()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
