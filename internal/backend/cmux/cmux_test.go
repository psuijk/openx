package cmux

import (
	"strings"
	"testing"

	"github.com/psuijk/openx/internal/config"
)

func TestBuild_NewWindow_BasicPlan(t *testing.T) {
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

	// First step creates workspace, then rename first tab, then new-surface + rename for second tab + send command
	// workspace(1) + rename(1) + new-surface(1) + rename(1) + send(1) = 5
	if len(p.Backend) != 5 {
		t.Errorf("Backend: got %d steps, want 5", len(p.Backend))
	}
	if !strings.Contains(p.Backend[0].Command, "cmux new-workspace") {
		t.Errorf("Backend[0]: expected new-workspace, got %q", p.Backend[0].Command)
	}
	if !strings.Contains(p.Backend[0].Command, "--cwd") {
		t.Errorf("Backend[0]: expected --cwd flag, got %q", p.Backend[0].Command)
	}

	if len(p.PostOpen) != 1 {
		t.Errorf("PostOpen: got %d steps, want 1", len(p.PostOpen))
	}
	if p.PostOpen[0].Command != "code ." {
		t.Errorf("PostOpen[0].Command: got %q, want %q", p.PostOpen[0].Command, "code .")
	}
}

func TestBuild_NewWindow_FirstTabCommand(t *testing.T) {
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

	// First tab command should be passed to new-workspace via --command
	if !strings.Contains(p.Backend[0].Command, "--command") {
		t.Errorf("expected --command in workspace creation, got %q", p.Backend[0].Command)
	}
}

func TestBuild_NewWindow_FirstTabNoCommand(t *testing.T) {
	b := &CmuxBackend{}
	cfg := config.Config{
		Name: "myproject",
		Path: "/some/path",
		Tabs: []config.Tab{
			{Name: "shell", Command: ""},
		},
	}

	p, err := b.Build(cfg, "new_window")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No command for first tab, so --command should not appear
	if strings.Contains(p.Backend[0].Command, "--command") {
		t.Errorf("expected no --command in workspace creation, got %q", p.Backend[0].Command)
	}
}

func TestBuild_NewWindow_NoTabs(t *testing.T) {
	b := &CmuxBackend{}
	cfg := config.Config{
		Name: "empty",
		Path: "/some/path",
	}

	p, err := b.Build(cfg, "new_window")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(p.Backend) != 1 {
		t.Errorf("Backend: got %d steps, want 1", len(p.Backend))
	}
}

func TestBuild_JoinMode(t *testing.T) {
	b := &CmuxBackend{}
	cfg := config.Config{
		Name: "testproject",
		Path: "/some/path",
		Tabs: []config.Tab{
			{Name: "shell", Command: ""},
			{Name: "claude", Command: "claude"},
		},
	}

	p, err := b.Build(cfg, "join")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Join mode should not create a workspace
	for _, step := range p.Backend {
		if strings.Contains(step.Command, "new-workspace") {
			t.Error("join mode should not create a new workspace")
		}
	}

	// Should have new-surface + rename for each tab, plus send for claude
	// shell: new-surface(1) + rename(1) = 2
	// claude: new-surface(1) + rename(1) + send(1) = 3
	// total = 5
	if len(p.Backend) != 5 {
		t.Errorf("Backend: got %d steps, want 5", len(p.Backend))
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
