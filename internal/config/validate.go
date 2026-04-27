package config

import (
	"errors"
	"fmt"
	"os"
)

// Validate checks a Config for required fields, valid values, and duplicate tab names.
func Validate(cfg Config) error {
	if cfg.Name == "" {
		return errors.New("config name cannot be empty")
	}

	if cfg.Path == "" {
		return errors.New("config path cannot be empty")
	}

	if _, err := os.Stat(cfg.Path); err != nil {
		return fmt.Errorf("invalid config path: %w", err)
	}

	if cfg.DefaultMode != "" && cfg.DefaultMode != "new_window" && cfg.DefaultMode != "join" {
		return errors.New("default mode must be either 'new_window' or 'join'")
	}

	seen := make(map[string]bool)
	for _, tab := range cfg.Tabs {
		if tab.Name == "" {
			return errors.New("tab name cannot be empty")
		}

		if seen[tab.Name] {
			return fmt.Errorf("duplicate tab name: %s", tab.Name)
		}

		seen[tab.Name] = true
	}

	return nil
}
