package command

import (
	"testing"

	"github.com/psuijk/openx/internal/config"
)

func TestClone_CopiesConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg := config.Config{
		Name:        "original",
		Path:        tmpDir,
		Backend:     "cmux",
		DefaultMode: "new_window",
		PreOpen:     []string{"git fetch"},
		PostOpen:    []string{"code ."},
		Tabs: []config.Tab{
			{Name: "claude", Command: "claude"},
			{Name: "shell", Command: ""},
		},
	}
	config.Store(cfg)

	err := cloneHandler([]string{"--path", tmpDir, "original", "cloned"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, err := config.Load("cloned")
	if err != nil {
		t.Fatalf("failed to load cloned config: %v", err)
	}
	if loaded.Name != "cloned" {
		t.Errorf("Name: got %q, want %q", loaded.Name, "cloned")
	}
	if loaded.Path != tmpDir {
		t.Errorf("Path: got %q, want %q", loaded.Path, tmpDir)
	}
	if len(loaded.Tabs) != 2 {
		t.Errorf("Tabs: got %d, want 2", len(loaded.Tabs))
	}
	if loaded.Tabs[0].Name != "claude" {
		t.Errorf("first tab: got %q, want %q", loaded.Tabs[0].Name, "claude")
	}
	if len(loaded.PreOpen) != 1 || loaded.PreOpen[0] != "git fetch" {
		t.Errorf("PreOpen: got %v", loaded.PreOpen)
	}
	if len(loaded.PostOpen) != 1 || loaded.PostOpen[0] != "code ." {
		t.Errorf("PostOpen: got %v", loaded.PostOpen)
	}
}

func TestClone_DefaultsPathToCwd(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg := config.Config{Name: "source", Path: tmpDir}
	config.Store(cfg)

	err := cloneHandler([]string{"source", "dest"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, _ := config.Load("dest")
	if loaded.Name != "dest" {
		t.Errorf("Name: got %q, want %q", loaded.Name, "dest")
	}
}

func TestClone_SourceNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	err := cloneHandler([]string{"nonexistent", "new"})
	if err == nil {
		t.Fatal("expected error for nonexistent source")
	}
}

func TestClone_NoArgs(t *testing.T) {
	err := cloneHandler([]string{})
	if err == nil {
		t.Fatal("expected error for no args")
	}
}

func TestClone_OneArg(t *testing.T) {
	err := cloneHandler([]string{"only-one"})
	if err == nil {
		t.Fatal("expected error for only one arg")
	}
}
