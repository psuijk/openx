package command

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/psuijk/openx/internal/config"
)

func removeHandler(args []string) error {
	fs := flag.NewFlagSet("remove", flag.ContinueOnError)

	confirm := fs.Bool("yes", false, "skip confirmation")
	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("parse remove flags: %w", err)
	}

	if len(fs.Args()) != 1 {
		return errors.New("usage: openx remove <project-name> [--yes]")
	}

	path, err := config.GetProjectConfigPath(fs.Args()[0])
	if err != nil {
		return fmt.Errorf("getting project config path: %w", err)
	}

	if !*confirm {
		fmt.Printf("Are you sure you want to remove %s? [y/n]: ", fs.Args()[0])
		var answer string
		fmt.Scanln(&answer)
		if answer != "yes" && answer != "y" {
			return nil
		}
	}

	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("removing project config: %w", err)
	}

	return nil
}
