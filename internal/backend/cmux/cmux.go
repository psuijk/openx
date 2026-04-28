package cmux

import (
	"fmt"

	"github.com/psuijk/openx/internal/backend"
	"github.com/psuijk/openx/internal/config"
	"github.com/psuijk/openx/internal/plan"
)

// CmuxBackend implements the Backend interface for cmux workspaces.
type CmuxBackend struct {
}

func init() {
	backend.Register("cmux", &CmuxBackend{})
}

func (cmx *CmuxBackend) Build(cfg config.Config, mode string) (*plan.Plan, error) {

	var preOpen []plan.Step
	var backendStps []plan.Step
	var postOpen []plan.Step
	for _, stp := range cfg.PreOpen {
		preOpen = append(preOpen, plan.Step{Command: stp, Description: stp})
	}

	backendStps = append(backendStps, plan.Step{Command: fmt.Sprintf("cmux new-workspace --name %s", cfg.Name), Description: fmt.Sprintf("create workspace %q",
		cfg.Name)})
	for _, tab := range cfg.Tabs {
		cmdStr := ""
		if tab.Command != "" {
			cmdStr = fmt.Sprintf(" --command %s", tab.Command)
		}

		backendStps = append(backendStps, plan.Step{Command: fmt.Sprintf("cmux new-surface --name %s%s", tab.Name, cmdStr), Description: fmt.Sprintf("create tab %q", tab.Name)})
	}

	for _, stp := range cfg.PostOpen {
		postOpen = append(postOpen, plan.Step{Command: stp, Description: stp})
	}

	return &plan.Plan{PreOpen: preOpen, Backend: backendStps, PostOpen: postOpen}, nil
}

func (cmx *CmuxBackend) PrintPlan(p *plan.Plan) error {
	if len(p.PreOpen) > 0 {
		fmt.Println("[pre_open]")
		for _, step := range p.PreOpen {
			fmt.Printf("  %-30s  %s\n", step.Description, step.Command)
		}
		fmt.Println()
	}

	if len(p.Backend) > 0 {
		fmt.Println("[backend]")
		for _, step := range p.Backend {
			fmt.Printf("  %-30s  %s\n", step.Description, step.Command)
		}
		fmt.Println()
	}

	if len(p.PostOpen) > 0 {
		fmt.Println("[post_open]")
		for _, step := range p.PostOpen {
			fmt.Printf("  %-30s  %s\n", step.Description, step.Command)
		}
		fmt.Println()
	}

	return nil
}
