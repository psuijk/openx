package cmux

import (
	"fmt"
	"os"
	"strings"

	"github.com/psuijk/openx/internal/backend"
	"github.com/psuijk/openx/internal/config"
	"github.com/psuijk/openx/internal/plan"
	"github.com/psuijk/openx/internal/shell"
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

// Execute overrides Base.Execute to handle cmux workspace/surface ID passing between steps.
func (cmx *CmuxBackend) Execute(p *plan.Plan, dir string) error {
	// Run pre_open steps (abort on failure)
	for _, stp := range p.PreOpen {
		err := shell.Execute(stp.Command, dir)
		if err != nil {
			return fmt.Errorf("pre_open failed: %w", err)
		}
	}

	// Run backend steps, capturing workspace/surface refs to pass between commands
	var workspaceRef string
	var surfaceRef string

	for _, stp := range p.Backend {
		cmd := stp.Command

		switch {
		case strings.HasPrefix(cmd, "cmux new-workspace"):
			out, err := shell.ExecuteCapture(cmd, dir)
			if err != nil {
				return fmt.Errorf("backend step failed: %w", err)
			}
			// Output: "OK workspace:11"
			if strings.HasPrefix(out, "OK ") {
				workspaceRef = strings.TrimPrefix(out, "OK ")
			}
			// The first surface in the new workspace — get it
			if workspaceRef != "" {
				surfaces, err := shell.ExecuteCapture(
					fmt.Sprintf("cmux list-pane-surfaces --workspace %s", workspaceRef), dir)
				if err == nil && surfaces != "" {
					// First line has the surface ref, e.g. "* surface:5  ..."
					firstLine := strings.Split(surfaces, "\n")[0]
					fields := strings.Fields(firstLine)
					for _, f := range fields {
						if strings.HasPrefix(f, "surface:") {
							surfaceRef = f
							break
						}
					}
				}
			}

		case strings.HasPrefix(cmd, "cmux new-surface"):
			nsCmd := cmd
			if workspaceRef != "" {
				nsCmd = fmt.Sprintf("cmux new-surface --workspace %s", workspaceRef)
			}
			out, err := shell.ExecuteCapture(nsCmd, dir)
			if err != nil {
				return fmt.Errorf("backend step failed: %w", err)
			}
			// Output: "OK surface:24 pane:14 workspace:15" — extract just the surface ref
			for _, f := range strings.Fields(out) {
				if strings.HasPrefix(f, "surface:") {
					surfaceRef = f
					break
				}
			}

		case strings.HasPrefix(cmd, "cmux rename-tab"):
			// Extract the title from the original command
			title := strings.TrimPrefix(cmd, "cmux rename-tab ")
			rtCmd := fmt.Sprintf("cmux rename-tab --surface %s", surfaceRef)
			if workspaceRef != "" {
				rtCmd += fmt.Sprintf(" --workspace %s", workspaceRef)
			}
			rtCmd += " " + title
			err := shell.Execute(rtCmd, dir)
			if err != nil {
				return fmt.Errorf("backend step failed: %w", err)
			}

		case strings.HasPrefix(cmd, "cmux send"):
			// Extract the text from the original command
			text := strings.TrimPrefix(cmd, "cmux send ")
			sCmd := fmt.Sprintf("cmux send --surface %s", surfaceRef)
			if workspaceRef != "" {
				sCmd += fmt.Sprintf(" --workspace %s", workspaceRef)
			}
			sCmd += " " + text
			err := shell.Execute(sCmd, dir)
			if err != nil {
				return fmt.Errorf("backend step failed: %w", err)
			}

		default:
			err := shell.Execute(cmd, dir)
			if err != nil {
				return fmt.Errorf("backend step failed: %w", err)
			}
		}
	}

	// Run post_open steps (log failures but continue)
	for _, stp := range p.PostOpen {
		err := shell.Execute(stp.Command, dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "post_open warning: %s\n", err)
		}
	}

	return nil
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
