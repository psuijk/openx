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

	// new-window(1) + rename-workspace(1) + rename-tab(1) + new-surface(1) + rename-tab(1) + send(1) = 6
	if len(p.Backend) != 6 {
		t.Errorf("Backend: got %d steps, want 6", len(p.Backend))
	}
	if p.Backend[0].Command != "cmux new-window" {
		t.Errorf("Backend[0]: expected new-window, got %q", p.Backend[0].Command)
	}
	if !strings.Contains(p.Backend[1].Command, "cmux rename-workspace") {
		t.Errorf("Backend[1]: expected rename-workspace, got %q", p.Backend[1].Command)
	}

	if len(p.PostOpen) != 1 {
		t.Errorf("PostOpen: got %d steps, want 1", len(p.PostOpen))
	}
}

func TestBuild_Default_NoNewWindow(t *testing.T) {
	b := &CmuxBackend{}
	cfg := config.Config{
		Name: "testproject",
		Path: "/some/path",
		Tabs: []config.Tab{
			{Name: "shell", Command: ""},
		},
	}

	p, err := b.Build(cfg, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, step := range p.Backend {
		if step.Command == "cmux new-window" {
			t.Error("default mode should not create a new window")
		}
	}
	if !strings.Contains(p.Backend[0].Command, "cmux new-workspace") {
		t.Errorf("Backend[0]: expected new-workspace, got %q", p.Backend[0].Command)
	}
}

func TestBuild_NewWindow_TabWithCommand(t *testing.T) {
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

	// new-window(1) + rename-workspace(1) + rename-tab(1) + send(1) = 4
	if len(p.Backend) != 4 {
		t.Errorf("Backend: got %d steps, want 4", len(p.Backend))
	}
	hasSend := false
	for _, step := range p.Backend {
		if strings.Contains(step.Command, "cmux send") {
			hasSend = true
		}
	}
	if !hasSend {
		t.Error("expected a send command for tab with command")
	}
}

func TestBuild_NewWindow_TabNoCommand(t *testing.T) {
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

	// new-window(1) + rename-workspace(1) + rename-tab(1) = 3
	if len(p.Backend) != 3 {
		t.Errorf("Backend: got %d steps, want 3", len(p.Backend))
	}
	for _, step := range p.Backend {
		if strings.Contains(step.Command, "cmux send") {
			t.Errorf("expected no send command, got %q", step.Command)
		}
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

	// new-window(1) + rename-workspace(1) = 2
	if len(p.Backend) != 2 {
		t.Errorf("Backend: got %d steps, want 2", len(p.Backend))
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

	for _, step := range p.Backend {
		if strings.Contains(step.Command, "new-workspace") {
			t.Error("join mode should not create a new workspace")
		}
		if strings.Contains(step.Command, "new-window") {
			t.Error("join mode should not create a new window")
		}
	}

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

	if p.Backend[0].Description != "create new cmux window" {
		t.Errorf("window description: got %q", p.Backend[0].Description)
	}
	if p.Backend[1].Description != `rename workspace to "myproject"` {
		t.Errorf("workspace description: got %q", p.Backend[1].Description)
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
