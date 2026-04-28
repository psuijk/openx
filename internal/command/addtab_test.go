package command

import (
	"testing"

	"github.com/psuijk/openx/internal/config"
)

func setupProject(t *testing.T, name string, tabs []config.Tab) {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)
	cfg := config.Config{
		Name: name,
		Path: tmpDir,
		Tabs: tabs,
	}
	if err := config.Store(cfg); err != nil {
		t.Fatalf("failed to store config: %v", err)
	}
}

func TestAddTab_NewTab(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)
	cfg := config.Config{Name: "proj", Path: tmpDir}
	config.Store(cfg)

	err := addTabHandler([]string{"proj", "shell"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, _ := config.Load("proj")
	if len(loaded.Tabs) != 1 {
		t.Fatalf("expected 1 tab, got %d", len(loaded.Tabs))
	}
	if loaded.Tabs[0].Name != "shell" {
		t.Errorf("tab name: got %q, want %q", loaded.Tabs[0].Name, "shell")
	}
}

func TestAddTab_WithCommand(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)
	cfg := config.Config{Name: "proj", Path: tmpDir}
	config.Store(cfg)

	err := addTabHandler([]string{"--command", "claude", "proj", "claude"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, _ := config.Load("proj")
	if loaded.Tabs[0].Command != "claude" {
		t.Errorf("tab command: got %q, want %q", loaded.Tabs[0].Command, "claude")
	}
}

func TestAddTab_UpdateExisting(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)
	cfg := config.Config{
		Name: "proj",
		Path: tmpDir,
		Tabs: []config.Tab{{Name: "claude", Command: "claude"}},
	}
	config.Store(cfg)

	err := addTabHandler([]string{"--command", "claude --model opus", "proj", "claude"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, _ := config.Load("proj")
	if len(loaded.Tabs) != 1 {
		t.Fatalf("expected 1 tab (updated, not duplicated), got %d", len(loaded.Tabs))
	}
	if loaded.Tabs[0].Command != "claude --model opus" {
		t.Errorf("tab command: got %q, want %q", loaded.Tabs[0].Command, "claude --model opus")
	}
}

func TestAddTab_AfterTab(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)
	cfg := config.Config{
		Name: "proj",
		Path: tmpDir,
		Tabs: []config.Tab{
			{Name: "first", Command: ""},
			{Name: "last", Command: ""},
		},
	}
	config.Store(cfg)

	err := addTabHandler([]string{"--after", "first", "proj", "middle"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, _ := config.Load("proj")
	if len(loaded.Tabs) != 3 {
		t.Fatalf("expected 3 tabs, got %d", len(loaded.Tabs))
	}
	if loaded.Tabs[1].Name != "middle" {
		t.Errorf("tab at index 1: got %q, want %q", loaded.Tabs[1].Name, "middle")
	}
}

func TestAddTab_BeforeTab(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)
	cfg := config.Config{
		Name: "proj",
		Path: tmpDir,
		Tabs: []config.Tab{
			{Name: "first", Command: ""},
			{Name: "last", Command: ""},
		},
	}
	config.Store(cfg)

	err := addTabHandler([]string{"--before", "last", "proj", "middle"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, _ := config.Load("proj")
	if len(loaded.Tabs) != 3 {
		t.Fatalf("expected 3 tabs, got %d", len(loaded.Tabs))
	}
	if loaded.Tabs[1].Name != "middle" {
		t.Errorf("tab at index 1: got %q, want %q", loaded.Tabs[1].Name, "middle")
	}
}

func TestAddTab_AfterNonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)
	cfg := config.Config{Name: "proj", Path: tmpDir}
	config.Store(cfg)

	err := addTabHandler([]string{"--after", "nope", "proj", "newtab"})
	if err == nil {
		t.Fatal("expected error for nonexistent --after tab")
	}
}

func TestAddTab_NoArgs(t *testing.T) {
	err := addTabHandler([]string{})
	if err == nil {
		t.Fatal("expected error for no args")
	}
}

func TestAddTab_BothAfterAndBefore(t *testing.T) {
	err := addTabHandler([]string{"--after", "a", "--before", "b", "proj", "tab"})
	if err == nil {
		t.Fatal("expected error when both --after and --before are set")
	}
}
