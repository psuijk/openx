package command

import (
	"os"
	"testing"

	"github.com/psuijk/openx/internal/config"
)

func TestRemoveHandler_WithYesFlag(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg := config.Config{Name: "removeme", Path: "/some/path"}
	if err := config.Store(cfg); err != nil {
		t.Fatalf("failed to store config: %v", err)
	}

	err := removeHandler([]string{"--yes", "removeme"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path, _ := config.GetProjectConfigPath("removeme")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("expected config file to be deleted")
	}
}

func TestRemoveHandler_NoArgs(t *testing.T) {
	err := removeHandler([]string{})
	if err == nil {
		t.Fatal("expected error when no project name given")
	}
}

func TestRemoveHandler_NonexistentProject(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create the projects dir so path resolution works
	cfg := config.Config{Name: "dummy", Path: "/tmp"}
	config.Store(cfg)

	err := removeHandler([]string{"--yes", "doesnotexist"})
	if err == nil {
		t.Fatal("expected error removing nonexistent project")
	}
}
