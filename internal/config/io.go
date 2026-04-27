package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Store writes a project config to disk as TOML, creating directories as needed.
func Store(cfg Config) error {
	projectsDir, err := GetProjectsDir()
	if err != nil {
		return fmt.Errorf("resolve projects directory: %w", err)
	}
	err = os.MkdirAll(projectsDir, 0755)
	if err != nil {
		return fmt.Errorf("create projects directory: %w", err)
	}

	path, err := GetProjectConfigPath(cfg.Name)
	if err != nil {
		return fmt.Errorf("resolve config path: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create config file: %w", err)
	}
	defer file.Close()

	err = toml.NewEncoder(file).Encode(cfg)
	if err != nil {
		return fmt.Errorf("write config to %s: %w", path, err)
	}

	return nil
}

func Load(name string) (*Config, error) {
	path, err := GetProjectConfigPath(name)
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer file.Close()

	var cfg Config

	_, err = toml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("read config from %s: %w", path, err)
	}

	return &cfg, nil
}
