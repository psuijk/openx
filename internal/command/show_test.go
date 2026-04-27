package command

import (
	"testing"

	"github.com/psuijk/openx/internal/config"
)

func TestShowHandler_PrintsConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg := config.Config{
		Name:    "showtest",
		Path:    "/home/user/showtest",
		Backend: "cmux",
		Tabs: []config.Tab{
			{Name: "shell", Command: ""},
		},
	}
	if err := config.Store(cfg); err != nil {
		t.Fatalf("failed to store config: %v", err)
	}

	err := showHandler([]string{"showtest"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShowHandler_NoArgs(t *testing.T) {
	err := showHandler([]string{})
	if err == nil {
		t.Fatal("expected error when no project name given")
	}
}

func TestShowHandler_TooManyArgs(t *testing.T) {
	err := showHandler([]string{"one", "two"})
	if err == nil {
		t.Fatal("expected error when too many args given")
	}
}

func TestShowHandler_NonexistentProject(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	err := showHandler([]string{"doesnotexist"})
	if err == nil {
		t.Fatal("expected error for nonexistent project")
	}
}
