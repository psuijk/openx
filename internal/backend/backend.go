package backend

import (
	"fmt"

	"github.com/psuijk/openx/internal/config"
	"github.com/psuijk/openx/internal/plan"
)

// Backend is the interface that terminal backends (cmux, tmux, etc.) must implement.
type Backend interface {
	Build(cfg config.Config, mode string) (*plan.Plan, error)
	PrintPlan(pln *plan.Plan) error
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
