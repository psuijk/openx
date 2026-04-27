package config

import (
	"os"
	"testing"
)

func validConfig(t *testing.T) Config {
	t.Helper()
	dir := t.TempDir()
	return Config{
		Name:        "testproject",
		Path:        dir,
		DefaultMode: "new_window",
		Backend:     "cmux",
		Tabs: []Tab{
			{Name: "shell", Command: ""},
			{Name: "claude", Command: "claude"},
		},
	}
}

func TestValidate_ValidConfig(t *testing.T) {
	cfg := validConfig(t)
	if err := Validate(cfg); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_EmptyName(t *testing.T) {
	cfg := validConfig(t)
	cfg.Name = ""
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestValidate_EmptyPath(t *testing.T) {
	cfg := validConfig(t)
	cfg.Path = ""
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestValidate_NonexistentPath(t *testing.T) {
	cfg := validConfig(t)
	cfg.Path = "/nonexistent/path/that/does/not/exist"
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for nonexistent path")
	}
}

func TestValidate_PathIsFile(t *testing.T) {
	cfg := validConfig(t)
	f, err := os.CreateTemp("", "validate-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()
	cfg.Path = f.Name()
	// Path exists but is a file — Validate currently allows this.
	// This test documents current behavior.
	if err := Validate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_InvalidDefaultMode(t *testing.T) {
	cfg := validConfig(t)
	cfg.DefaultMode = "invalid"
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for invalid default mode")
	}
}

func TestValidate_EmptyDefaultMode(t *testing.T) {
	cfg := validConfig(t)
	cfg.DefaultMode = ""
	if err := Validate(cfg); err != nil {
		t.Fatalf("empty default mode should be valid, got: %v", err)
	}
}

func TestValidate_JoinMode(t *testing.T) {
	cfg := validConfig(t)
	cfg.DefaultMode = "join"
	if err := Validate(cfg); err != nil {
		t.Fatalf("join mode should be valid, got: %v", err)
	}
}

func TestValidate_DuplicateTabNames(t *testing.T) {
	cfg := validConfig(t)
	cfg.Tabs = []Tab{
		{Name: "shell", Command: ""},
		{Name: "shell", Command: "bash"},
	}
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for duplicate tab names")
	}
}

func TestValidate_EmptyTabName(t *testing.T) {
	cfg := validConfig(t)
	cfg.Tabs = []Tab{
		{Name: "", Command: "bash"},
	}
	if err := Validate(cfg); err == nil {
		t.Fatal("expected error for empty tab name")
	}
}

func TestValidate_NoTabs(t *testing.T) {
	cfg := validConfig(t)
	cfg.Tabs = nil
	if err := Validate(cfg); err != nil {
		t.Fatalf("config with no tabs should be valid, got: %v", err)
	}
}
