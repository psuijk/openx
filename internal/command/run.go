package command

import (
	"errors"
	"flag"
	"fmt"

	"github.com/psuijk/openx/internal/backend"
	"github.com/psuijk/openx/internal/config"
)

func runHandler(args []string) error {
	// parse flags (--dry-run, --join, --new-window, --backend)
	fs := flag.NewFlagSet("run", flag.ContinueOnError)

	dryRun := fs.Bool("dry-run", false, "print what would happen without executing")
	join := fs.Bool("join", false, "attach to current window instead of opening a new one")
	newWindow := fs.Bool("new-window", false, "force opening a new window")
	backendFlag := fs.String("backend", "", "backend to use (overrides project and global config)")
	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("parse run flags: %w", err)
	}

	if len(fs.Args()) == 0 {
		return errors.New("usage: openx <project-name> [--dry-run] [--join] [--new-window] [--backend NAME]")
	}

	cfg, err := config.Load(fs.Args()[0])
	if err != nil {
		return fmt.Errorf("load project config %q: %w", fs.Args()[0], err)
	}

	glblCfg, err := config.LoadGlobal()
	if err != nil {
		return fmt.Errorf("load global config: %w", err)
	}

	var backendName string
	if *backendFlag != "" {
		backendName = *backendFlag
	} else if cfg.Backend != "" {
		backendName = cfg.Backend
	} else if glblCfg.DefaultBackend != "" {
		backendName = glblCfg.DefaultBackend
	} else {
		backendName = "cmux"
	}
	actBackend, err := backend.Get(backendName)
	if err != nil {
		return fmt.Errorf("unknown backend %q: %w", backendName, err)
	}

	if *join && *newWindow {
		return errors.New("cannot use --join and --new-window together")
	}

	var mode string
	if *join {
		mode = "join"
	} else if *newWindow {
		mode = "new_window"
	} else if cfg.DefaultMode != "" {
		mode = cfg.DefaultMode
	} else if glblCfg.DefaultMode != "" {
		mode = glblCfg.DefaultMode
	} else {
		mode = "new_window"
	}

	p, err := actBackend.Build(*cfg, mode)
	if err != nil {
		return fmt.Errorf("building plan: %w", err)
	}

	if *dryRun {
		err = actBackend.PrintPlan(p)
		if err != nil {
			return fmt.Errorf("printing plan: %w", err)
		}
	} else {
		err = actBackend.Execute(p, cfg.Path)
		if err != nil {
			return fmt.Errorf("executing plan: %w", err)
		}
	}

	return nil
}
