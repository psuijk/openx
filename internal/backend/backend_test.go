package backend

import (
	"testing"

	"github.com/psuijk/openx/internal/config"
	"github.com/psuijk/openx/internal/plan"
)

// fakeBackend satisfies the Backend interface for testing.
type fakeBackend struct {
	Base
}

func (f *fakeBackend) Build(cfg config.Config, mode string) (*plan.Plan, error) {
	return &plan.Plan{}, nil
}

func (f *fakeBackend) PrintPlan(p *plan.Plan) error {
	return nil
}

func TestExecute_RunsStepsInOrder(t *testing.T) {
	dir := t.TempDir()
	p := &plan.Plan{
		PreOpen: []plan.Step{
			{Command: "echo pre", Description: "pre"},
		},
		Backend: []plan.Step{
			{Command: "echo backend", Description: "backend"},
		},
		PostOpen: []plan.Step{
			{Command: "echo post", Description: "post"},
		},
	}

	b := &Base{}
	err := b.Execute(p, dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_PreOpenFailureAborts(t *testing.T) {
	dir := t.TempDir()
	p := &plan.Plan{
		PreOpen: []plan.Step{
			{Command: "false", Description: "will fail"},
		},
		Backend: []plan.Step{
			{Command: "echo should not run", Description: "backend"},
		},
	}

	b := &Base{}
	err := b.Execute(p, dir)
	if err == nil {
		t.Fatal("expected error from failing pre_open step")
	}
}

func TestExecute_PostOpenFailureContinues(t *testing.T) {
	dir := t.TempDir()
	p := &plan.Plan{
		PostOpen: []plan.Step{
			{Command: "false", Description: "will fail"},
			{Command: "echo still runs", Description: "second"},
		},
	}

	b := &Base{}
	err := b.Execute(p, dir)
	if err != nil {
		t.Fatalf("post_open failure should not return error, got: %v", err)
	}
}

func TestRegisterAndGet(t *testing.T) {
	Register("fake", &fakeBackend{})

	b, err := Get("fake")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected backend, got nil")
	}
}

func TestGet_Unknown(t *testing.T) {
	_, err := Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown backend")
	}
}
