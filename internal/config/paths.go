package config

import (
	"fmt"
	"os"
	"path/filepath"
)

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

func GetProjectsDir() (string, error) {
	configPath, err := GetBaseConfigDir()
	if err != nil {
		return "", fmt.Errorf("get projects directory: %w", err)
	}

	return filepath.Join(configPath, "projects"), nil
}

func GetProjectConfigPath(projectName string) (string, error) {
	projectsDir, err := GetProjectsDir()
	if err != nil {
		return "", fmt.Errorf("get project config path: %w", err)
	}

	return filepath.Join(projectsDir, projectName+".toml"), nil
}
