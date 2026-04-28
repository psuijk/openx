package cmux

import (
	"testing"

	"github.com/psuijk/openx/internal/config"
)

func TestBuild_BasicPlan(t *testing.T) {
	b := &CmuxBackend{}
	cfg := config.Config{
		Name: "testproject",
		Path: "/some/path",
		Tabs: []config.Tab{
			{Name: "shell", Command: ""},
			{Name: "claude", Command: "claude"},
		},
		PreOpen:  []string{"git fetch"},
		PostOpen: []string{"code ."},
	}

	p, err := b.Build(cfg, "new_window")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(p.PreOpen) != 1 {
		t.Errorf("PreOpen: got %d steps, want 1", len(p.PreOpen))
	}
	if p.PreOpen[0].Command != "git fetch" {
		t.Errorf("PreOpen[0].Command: got %q, want %q", p.PreOpen[0].Command, "git fetch")
	}

	// 1 workspace + 2 tabs = 3 backend steps
	if len(p.Backend) != 3 {
		t.Errorf("Backend: got %d steps, want 3", len(p.Backend))
	}
	if p.Backend[0].Command != "cmux new-workspace --name testproject" {
		t.Errorf("Backend[0].Command: got %q, want workspace creation command", p.Backend[0].Command)
	}
	if p.Backend[1].Command != "cmux new-surface --name shell" {
		t.Errorf("Backend[1].Command: got %q, want surface without --command", p.Backend[1].Command)
	}
	if p.Backend[2].Command != "cmux new-surface --name claude --command claude" {
		t.Errorf("Backend[2].Command: got %q, want surface with --command", p.Backend[2].Command)
	}

	if len(p.PostOpen) != 1 {
		t.Errorf("PostOpen: got %d steps, want 1", len(p.PostOpen))
	}
	if p.PostOpen[0].Command != "code ." {
		t.Errorf("PostOpen[0].Command: got %q, want %q", p.PostOpen[0].Command, "code .")
	}
}

func TestBuild_NoTabs(t *testing.T) {
	b := &CmuxBackend{}
	cfg := config.Config{
		Name: "empty",
		Path: "/some/path",
	}

	p, err := b.Build(cfg, "new_window")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Just the workspace creation step
	if len(p.Backend) != 1 {
		t.Errorf("Backend: got %d steps, want 1", len(p.Backend))
	}
}

func TestBuild_NoPrePostOpen(t *testing.T) {
	b := &CmuxBackend{}
	cfg := config.Config{
		Name: "minimal",
		Path: "/some/path",
		Tabs: []config.Tab{
			{Name: "shell", Command: ""},
		},
	}

	p, err := b.Build(cfg, "new_window")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(p.PreOpen) != 0 {
		t.Errorf("PreOpen: got %d steps, want 0", len(p.PreOpen))
	}
	if len(p.PostOpen) != 0 {
		t.Errorf("PostOpen: got %d steps, want 0", len(p.PostOpen))
	}
}

func TestBuild_Descriptions(t *testing.T) {
	b := &CmuxBackend{}
	cfg := config.Config{
		Name: "myproject",
		Path: "/some/path",
		Tabs: []config.Tab{
			{Name: "dev", Command: "npm start"},
		},
	}

	p, err := b.Build(cfg, "new_window")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Backend[0].Description != `create workspace "myproject"` {
		t.Errorf("workspace description: got %q", p.Backend[0].Description)
	}
	if p.Backend[1].Description != `create tab "dev"` {
		t.Errorf("tab description: got %q", p.Backend[1].Description)
	}
}

func TestPrintPlan_NoError(t *testing.T) {
	b := &CmuxBackend{}
	cfg := config.Config{
		Name:     "testproject",
		Path:     "/some/path",
		PreOpen:  []string{"git fetch"},
		PostOpen: []string{"code ."},
		Tabs: []config.Tab{
			{Name: "shell", Command: ""},
		},
	}

	p, _ := b.Build(cfg, "new_window")
	err := b.PrintPlan(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
