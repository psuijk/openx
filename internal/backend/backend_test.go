package backend

import (
	"testing"

	"github.com/psuijk/openx/internal/config"
	"github.com/psuijk/openx/internal/plan"
)

// fakeBackend satisfies the Backend interface for testing.
type fakeBackend struct{}

func (f *fakeBackend) Build(cfg config.Config, mode string) (*plan.Plan, error) {
	return &plan.Plan{}, nil
}

func (f *fakeBackend) PrintPlan(p *plan.Plan) error {
	return nil
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
