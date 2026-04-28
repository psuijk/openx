package shell

import (
	"fmt"
	"os"
	"os/exec"
)

// Execute runs a shell command in the given directory, with stdin/stdout/stderr connected to the terminal.
func Execute(command string, dir string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command %q failed: %w", command, err)
	}
	return nil
}
