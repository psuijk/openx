package command

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/psuijk/openx/internal/config"
)

func TestAddHandler_CreatesConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	projectDir := t.TempDir()
	err := addHandler([]string{"--path", projectDir, "testproject"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := filepath.Join(tmpDir, "openx", "projects", "testproject.toml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("expected config file at %s", path)
	}

	var cfg config.Config
	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		t.Fatalf("failed to decode config: %v", err)
	}
	if cfg.Name != "testproject" {
		t.Errorf("Name: got %q, want %q", cfg.Name, "testproject")
	}
	if cfg.Path != projectDir {
		t.Errorf("Path: got %q, want %q", cfg.Path, projectDir)
	}
}

func TestAddHandler_DefaultsPathToCwd(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	err := addHandler([]string{"cwdproject"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := filepath.Join(tmpDir, "openx", "projects", "cwdproject.toml")
	var cfg config.Config
	_, err = toml.DecodeFile(path, &cfg)
	if err != nil {
		t.Fatalf("failed to decode config: %v", err)
	}

	cwd, _ := os.Getwd()
	if cfg.Path != cwd {
		t.Errorf("Path: got %q, want cwd %q", cfg.Path, cwd)
	}
}

func TestAddHandler_NoArgs(t *testing.T) {
	err := addHandler([]string{})
	if err == nil {
		t.Fatal("expected error when no project name given")
	}
}

func TestAddHandler_TooManyArgs(t *testing.T) {
	err := addHandler([]string{"one", "two"})
	if err == nil {
		t.Fatal("expected error when too many args given")
	}
}
