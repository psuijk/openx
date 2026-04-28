package backend

import (
	"fmt"
	"os"

	"github.com/psuijk/openx/internal/config"
	"github.com/psuijk/openx/internal/plan"
	"github.com/psuijk/openx/internal/shell"
)

// Base provides shared execution logic that backends inherit via embedding.
type Base struct{}

// Backend is the interface that terminal backends (cmux, tmux, etc.) must implement.
type Backend interface {
	Build(cfg config.Config, mode string) (*plan.Plan, error)
	PrintPlan(pln *plan.Plan) error
	Execute(plan *plan.Plan, dir string) error
}

// Execute runs all plan steps in order. Pre_open failures abort; post_open failures are logged but non-fatal.
func (b *Base) Execute(plan *plan.Plan, dir string) error {
	for _, stp := range plan.PreOpen {
		err := shell.Execute(stp.Command, dir)
		if err != nil {
			return fmt.Errorf("pre_open failed: %w", err)
		}
	}

	for _, stp := range plan.Backend {
		err := shell.Execute(stp.Command, dir)
		if err != nil {
			return fmt.Errorf("backend step failed: %w", err)
		}
	}

	for _, stp := range plan.PostOpen {
		err := shell.Execute(stp.Command, dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "post_open warning: %s\n", err)
		}
	}

	return nil
}

var registry = map[string]Backend{}

// Register adds a backend to the registry under the given name.
func Register(name string, b Backend) {
	registry[name] = b
}

// Get returns the backend registered under the given name, or an error if not found.
func Get(name string) (Backend, error) {
	backend, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown backend: %s", name)
	}

	return backend, nil
}
