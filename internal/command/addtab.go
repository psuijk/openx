package command

import (
	"errors"
	"flag"
	"fmt"

	"github.com/psuijk/openx/internal/config"
)

func addTabHandler(args []string) error {
	fs := flag.NewFlagSet("add-tab", flag.ContinueOnError)
	command := fs.String("command", "", "command to run in the tab")
	after := fs.String("after", "", "insert after this tab")
	before := fs.String("before", "", "insert before this tab")
	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("parse add-tab flags: %w", err)
	}

	if len(fs.Args()) != 2 {
		return errors.New("usage: openx add-tab <project-name> <tab-name> [--command CMD] [--after TAB] [--before TAB]")
	}

	if *after != "" && *before != "" {
		return errors.New("cannot use --after and --before together")
	}

	projectName := fs.Args()[0]
	tabName := fs.Args()[1]

	cfg, err := config.Load(projectName)
	if err != nil {
		return fmt.Errorf("load project config %q: %w", projectName, err)
	}

	// Check if tab already exists — update in place
	for i, tab := range cfg.Tabs {
		if tab.Name == tabName {
			cfg.Tabs[i].Command = *command
			return config.Store(*cfg)
		}
	}

	// New tab
	newTab := config.Tab{Name: tabName, Command: *command}

	if *after != "" {
		inserted := false
		for i, tab := range cfg.Tabs {
			if tab.Name == *after {
				// Insert after position i
				cfg.Tabs = append(cfg.Tabs[:i+1], append([]config.Tab{newTab}, cfg.Tabs[i+1:]...)...)
				inserted = true
				break
			}
		}
		if !inserted {
			return fmt.Errorf("tab %q not found", *after)
		}
	} else if *before != "" {
		inserted := false
		for i, tab := range cfg.Tabs {
			if tab.Name == *before {
				// Insert before position i
				cfg.Tabs = append(cfg.Tabs[:i], append([]config.Tab{newTab}, cfg.Tabs[i:]...)...)
				inserted = true
				break
			}
		}
		if !inserted {
			return fmt.Errorf("tab %q not found", *before)
		}
	} else {
		cfg.Tabs = append(cfg.Tabs, newTab)
	}

	err = config.Store(*cfg)
	if err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	fmt.Printf("added tab %q to %s\n", tabName, projectName)
	return nil
}
