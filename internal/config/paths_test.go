package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetBaseConfigDir_WithXDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/custom-config")

	dir, err := GetBaseConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join("/tmp/custom-config", "openx")
	if dir != want {
		t.Errorf("got %q, want %q", dir, want)
	}
}

func TestGetBaseConfigDir_WithoutXDG(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")

	dir, err := GetBaseConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".config", "openx")
	if dir != want {
		t.Errorf("got %q, want %q", dir, want)
	}
}

func TestGetProjectsDir(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/custom-config")

	dir, err := GetProjectsDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join("/tmp/custom-config", "openx", "projects")
	if dir != want {
		t.Errorf("got %q, want %q", dir, want)
	}
}

func TestGetProjectConfigPath(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/custom-config")

	path, err := GetProjectConfigPath("myproject")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(path, "myproject.toml") {
		t.Errorf("expected path to end with myproject.toml, got %q", path)
	}
	want := filepath.Join("/tmp/custom-config", "openx", "projects", "myproject.toml")
	if path != want {
		t.Errorf("got %q, want %q", path, want)
	}
}
