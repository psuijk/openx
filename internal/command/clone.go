package command

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/psuijk/openx/internal/config"
)

func cloneHandler(args []string) error {
	fs := flag.NewFlagSet("clone", flag.ContinueOnError)
	path := fs.String("path", "", "project directory path for the new config")
	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("parse clone flags: %w", err)
	}

	if len(fs.Args()) != 2 {
		return errors.New("usage: openx clone <source-project> <new-name> [--path PATH]")
	}

	sourceName := fs.Args()[0]
	newName := fs.Args()[1]

	cfg, err := config.Load(sourceName)
	if err != nil {
		return fmt.Errorf("load source config %q: %w", sourceName, err)
	}

	cfg.Name = newName

	if *path != "" {
		cfg.Path = *path
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get current directory: %w", err)
		}
		cfg.Path = cwd
	}

	err = config.Store(*cfg)
	if err != nil {
		return fmt.Errorf("save cloned config: %w", err)
	}

	fmt.Printf("cloned %q to %q\n", sourceName, newName)
	return nil
}
