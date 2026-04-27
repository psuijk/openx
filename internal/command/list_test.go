package command

import (
	"os"
	"testing"

	"github.com/psuijk/openx/internal/config"
)

func TestListHandler_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Create the projects directory so listHandler doesn't error
	cfg := config.Config{Name: "dummy", Path: "/tmp"}
	config.Store(cfg)

	// Remove the file, leaving an empty directory
	path, _ := config.GetProjectConfigPath("dummy")
	os.Remove(path)

	err := listHandler()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListHandler_WithProjects(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	projects := []config.Config{
		{Name: "alpha", Path: "/path/alpha"},
		{Name: "beta", Path: "/path/beta"},
	}
	for _, p := range projects {
		if err := config.Store(p); err != nil {
			t.Fatalf("failed to store %q: %v", p.Name, err)
		}
	}

	err := listHandler()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListHandler_NoDirError(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Projects dir doesn't exist yet
	err := listHandler()
	if err == nil {
		t.Fatal("expected error when projects directory doesn't exist")
	}
}
