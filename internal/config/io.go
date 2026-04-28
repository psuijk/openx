package config

import (
	"fmt"
	"os"
	"path/filepath"

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

// Load reads a project's TOML config from disk and returns it.
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

// LoadGlobal reads the global config, returning sensible defaults if the file doesn't exist.
func LoadGlobal() (*GlobalConfig, error) {
	basePath, err := GetBaseConfigDir()
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}

	glblcfgPath := filepath.Join(basePath, "config.toml")

	file, err := os.Open(glblcfgPath)
	if os.IsNotExist(err) {
		return &GlobalConfig{DefaultBackend: "cmux", DefaultMode: "new_window"}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer file.Close()

	var glblcfg GlobalConfig
	_, err = toml.NewDecoder(file).Decode(&glblcfg)
	if err != nil {
		return nil, fmt.Errorf("read config from %s: %w", glblcfgPath, err)
	}

	return &glblcfg, nil
}
