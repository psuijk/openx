package command

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/psuijk/openx/internal/config"
)

func showHandler(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: openx show <project-name>")
	}

	cfg, err := config.Load(args[0])
	if err != nil {
		return fmt.Errorf("load project config %q: %w", args[0], err)
	}

	return toml.NewEncoder(os.Stdout).Encode(cfg)
}
