package shell

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExecute_SimpleCommand(t *testing.T) {
	dir := t.TempDir()
	err := Execute("echo hello", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_RunsInDir(t *testing.T) {
	dir := t.TempDir()
	err := Execute("pwd > output.txt", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "output.txt"))
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	// pwd output should match the dir we passed
	got := string(content)
	if got[:len(got)-1] != dir { // trim trailing newline
		t.Errorf("command ran in %q, want %q", got, dir)
	}
}

func TestExecute_FailingCommand(t *testing.T) {
	dir := t.TempDir()
	err := Execute("false", dir)
	if err == nil {
		t.Fatal("expected error for failing command")
	}
}

func TestExecute_ShellFeatures(t *testing.T) {
	dir := t.TempDir()
	// Test that shell features (&&, pipes) work
	err := Execute("echo hello && echo world", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_InvalidDir(t *testing.T) {
	err := Execute("echo hello", "/nonexistent/path")
	if err == nil {
		t.Fatal("expected error for invalid directory")
	}
}
