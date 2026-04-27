package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/psuijk/openx/internal/config"
)

func listHandler() error {
	path, err := config.GetProjectsDir()
	if err != nil {
		return fmt.Errorf("resolve projects directory: %w", err)
	}

	dir, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("read projects directory: %w", err)
	}

	for _, proj := range dir {
		cfg, err := config.Load(strings.TrimSuffix(proj.Name(), ".toml"))
		if err != nil {
			return fmt.Errorf("load project config %q: %w", proj.Name(), err)
		}
		fmt.Printf("%-20s %s\n", cfg.Name, cfg.Path)
	}

	return nil
}
