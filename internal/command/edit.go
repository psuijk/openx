package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/psuijk/openx/internal/config"
)

func editHandler(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: openx edit <project-name>")
	}

	path, err := config.GetProjectConfigPath(args[0])
	if err != nil {
		return fmt.Errorf("resolve config path: %w", err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		return errors.New("$EDITOR not set")
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("opening editor: %w", err)
	}

	return nil
}
