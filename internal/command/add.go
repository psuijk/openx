package command

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/psuijk/openx/internal/config"
)

func addHandler(args []string) error {

	fs := flag.NewFlagSet("add flags", flag.ContinueOnError)

	curDir, err := os.Getwd()

	if err != nil {
		return fmt.Errorf("get current directory: %w", err)
	}

	path := fs.String("path", curDir, "project directory path")

	err = fs.Parse(args)

	if err != nil {
		return fmt.Errorf("parse add flags: %w", err)
	}

	if len(fs.Args()) == 0 {
		return errors.New("usage: openx add <project-name> [--path PATH]")
	}

	cfg := config.Config{Name: fs.Args()[0], Path: *path}

	err = config.Store(cfg)

	if err != nil {
		return fmt.Errorf("saving new config: %w", err)
	}

	return nil
}
