package cmux

import (
	"fmt"

	"github.com/psuijk/openx/internal/backend"
	"github.com/psuijk/openx/internal/config"
	"github.com/psuijk/openx/internal/plan"
)

// CmuxBackend implements the Backend interface for cmux workspaces.
type CmuxBackend struct {
	backend.Base
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

	if mode == "join" {
		// Join mode: add tabs to the current workspace (CMUX_WORKSPACE_ID is auto-set)
		for _, tab := range cfg.Tabs {
			backendStps = append(backendStps, plan.Step{
				Command:     "cmux new-surface",
				Description: fmt.Sprintf("create tab %q", tab.Name),
			})
			backendStps = append(backendStps, plan.Step{
				Command:     fmt.Sprintf("cmux rename-tab %q", tab.Name),
				Description: fmt.Sprintf("rename tab to %q", tab.Name),
			})
			if tab.Command != "" {
				backendStps = append(backendStps, plan.Step{
					Command:     fmt.Sprintf("cmux send %q", tab.Command),
					Description: fmt.Sprintf("run %q in tab %q", tab.Command, tab.Name),
				})
			}
		}
	} else {
		// New window mode: create a workspace, then add tabs
		newWsCmd := fmt.Sprintf("cmux new-workspace --name %q --cwd %q", cfg.Name, cfg.Path)

		// First tab gets created with the workspace
		if len(cfg.Tabs) > 0 {
			first := cfg.Tabs[0]
			if first.Command != "" {
				newWsCmd += fmt.Sprintf(" --command %q", first.Command)
			}
			backendStps = append(backendStps, plan.Step{
				Command:     newWsCmd,
				Description: fmt.Sprintf("create workspace %q", cfg.Name),
			})
			backendStps = append(backendStps, plan.Step{
				Command:     fmt.Sprintf("cmux rename-tab %q", first.Name),
				Description: fmt.Sprintf("rename tab to %q", first.Name),
			})

			// Remaining tabs
			for _, tab := range cfg.Tabs[1:] {
				backendStps = append(backendStps, plan.Step{
					Command:     "cmux new-surface",
					Description: fmt.Sprintf("create tab %q", tab.Name),
				})
				backendStps = append(backendStps, plan.Step{
					Command:     fmt.Sprintf("cmux rename-tab %q", tab.Name),
					Description: fmt.Sprintf("rename tab to %q", tab.Name),
				})
				if tab.Command != "" {
					backendStps = append(backendStps, plan.Step{
						Command:     fmt.Sprintf("cmux send %q", tab.Command),
						Description: fmt.Sprintf("run %q in tab %q", tab.Command, tab.Name),
					})
				}
			}
		} else {
			// No tabs, just create the workspace
			backendStps = append(backendStps, plan.Step{
				Command:     newWsCmd,
				Description: fmt.Sprintf("create workspace %q", cfg.Name),
			})
		}
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
