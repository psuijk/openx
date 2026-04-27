package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestStore_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg := Config{
		Name: "testproject",
		Path: "/some/path",
	}

	err := Store(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := filepath.Join(tmpDir, "openx", "projects", "testproject.toml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("expected config file at %s, but it does not exist", path)
	}
}

func TestStore_WritesValidTOML(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg := Config{
		Name:        "myapp",
		Path:        "/home/user/myapp",
		DefaultMode: "new_window",
		Backend:     "cmux",
		PreOpen:     []string{"git fetch"},
		PostOpen:    []string{"code ."},
		Tabs: []Tab{
			{Name: "shell", Command: ""},
			{Name: "claude", Command: "claude"},
		},
	}

	err := Store(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := filepath.Join(tmpDir, "openx", "projects", "myapp.toml")
	var loaded Config
	_, err = toml.DecodeFile(path, &loaded)
	if err != nil {
		t.Fatalf("failed to decode written TOML: %v", err)
	}

	if loaded.Name != cfg.Name {
		t.Errorf("Name: got %q, want %q", loaded.Name, cfg.Name)
	}
	if loaded.Path != cfg.Path {
		t.Errorf("Path: got %q, want %q", loaded.Path, cfg.Path)
	}
	if loaded.DefaultMode != cfg.DefaultMode {
		t.Errorf("DefaultMode: got %q, want %q", loaded.DefaultMode, cfg.DefaultMode)
	}
	if loaded.Backend != cfg.Backend {
		t.Errorf("Backend: got %q, want %q", loaded.Backend, cfg.Backend)
	}
	if len(loaded.Tabs) != len(cfg.Tabs) {
		t.Errorf("Tabs count: got %d, want %d", len(loaded.Tabs), len(cfg.Tabs))
	}
	if len(loaded.PreOpen) != 1 || loaded.PreOpen[0] != "git fetch" {
		t.Errorf("PreOpen: got %v, want %v", loaded.PreOpen, cfg.PreOpen)
	}
	if len(loaded.PostOpen) != 1 || loaded.PostOpen[0] != "code ." {
		t.Errorf("PostOpen: got %v, want %v", loaded.PostOpen, cfg.PostOpen)
	}
}

func TestLoad_ReturnsConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	original := Config{
		Name:        "loadtest",
		Path:        "/home/user/loadtest",
		DefaultMode: "join",
		Backend:     "cmux",
		PreOpen:     []string{"git pull"},
		PostOpen:    []string{"code ."},
		Tabs: []Tab{
			{Name: "shell", Command: ""},
			{Name: "dev", Command: "npm start"},
		},
	}

	err := Store(original)
	if err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	loaded, err := Load("loadtest")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Name != original.Name {
		t.Errorf("Name: got %q, want %q", loaded.Name, original.Name)
	}
	if loaded.Path != original.Path {
		t.Errorf("Path: got %q, want %q", loaded.Path, original.Path)
	}
	if loaded.DefaultMode != original.DefaultMode {
		t.Errorf("DefaultMode: got %q, want %q", loaded.DefaultMode, original.DefaultMode)
	}
	if loaded.Backend != original.Backend {
		t.Errorf("Backend: got %q, want %q", loaded.Backend, original.Backend)
	}
	if len(loaded.Tabs) != len(original.Tabs) {
		t.Errorf("Tabs count: got %d, want %d", len(loaded.Tabs), len(original.Tabs))
	}
	if len(loaded.PreOpen) != 1 || loaded.PreOpen[0] != "git pull" {
		t.Errorf("PreOpen: got %v, want %v", loaded.PreOpen, original.PreOpen)
	}
	if len(loaded.PostOpen) != 1 || loaded.PostOpen[0] != "code ." {
		t.Errorf("PostOpen: got %v, want %v", loaded.PostOpen, original.PostOpen)
	}
}

func TestLoad_NonexistentProject(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	_, err := Load("doesnotexist")
	if err == nil {
		t.Fatal("expected error loading nonexistent project, got nil")
	}
}

func TestStore_CreatesDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	projectsDir := filepath.Join(tmpDir, "openx", "projects")
	if _, err := os.Stat(projectsDir); !os.IsNotExist(err) {
		t.Fatal("projects dir should not exist before Store")
	}

	cfg := Config{Name: "newproject", Path: "/some/path"}
	err := Store(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(projectsDir); os.IsNotExist(err) {
		t.Fatal("Store should have created the projects directory")
	}
}
