package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetBaseConfigDir returns the openx config directory, using $XDG_CONFIG_HOME or ~/.config/openx.
func GetBaseConfigDir() (string, error) {
	if basePath := os.Getenv("XDG_CONFIG_HOME"); basePath != "" {
		return filepath.Join(basePath, "openx"), nil
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get base config directory: %w", err)
	}

	return filepath.Join(userHome, ".config", "openx"), nil
}

// GetProjectsDir returns the directory where project config files are stored.
func GetProjectsDir() (string, error) {
	configPath, err := GetBaseConfigDir()
	if err != nil {
		return "", fmt.Errorf("get projects directory: %w", err)
	}

	return filepath.Join(configPath, "projects"), nil
}

// GetProjectConfigPath returns the file path for a project's TOML config.
func GetProjectConfigPath(projectName string) (string, error) {
	projectsDir, err := GetProjectsDir()
	if err != nil {
		return "", fmt.Errorf("get project config path: %w", err)
	}

	return filepath.Join(projectsDir, projectName+".toml"), nil
}
