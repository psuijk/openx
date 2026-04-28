# openx

A CLI tool for opening project-specific terminal workspaces. Run `openx <project-name>` to spin up your full dev environment — terminal tabs, editor, background tasks — all from a single command.

## Install

```bash
go install github.com/psuijk/openx/cmd/openx@latest
```

Make sure `~/go/bin` is in your PATH. If not, add this to your shell config (`~/.zshrc` or `~/.bashrc`):

```bash
export PATH="$HOME/go/bin:$PATH"
```

## Usage

```
openx add <name> [--path PATH]   # create a project config (defaults to cwd)
openx list                       # list all projects
openx show <name>                # print a project's config
openx edit <name>                # open config in $EDITOR
openx remove <name> [--yes]      # delete a project config
openx <name>                     # open the project workspace
openx <name> --dry-run           # print what would happen
openx <name> --join              # add tabs to current workspace
openx <name> --new-window        # force new workspace
openx <name> --backend <name>    # override backend
openx version                    # print version
```

Note: flags must come before the project name (e.g. `openx --dry-run myproject`).

## Config

Configs are TOML files stored at `$XDG_CONFIG_HOME/openx/projects/<name>.toml` (defaults to `~/.config/openx/projects/`).

```toml
name = "myproject"
path = "/Users/you/code/myproject"
default_mode = "new_window"
backend = "cmux"

pre_open = ["git fetch"]
post_open = ["code ."]

[[tabs]]
name = "claude"
command = "claude"

[[tabs]]
name = "backend"
command = "task start-backend"

[[tabs]]
name = "shell"
```

A global config at `~/.config/openx/config.toml` sets defaults for all projects:

```toml
default_mode = "new_window"
default_backend = "cmux"
```

## Requirements

- [Go](https://go.dev/) 1.26+
- [cmux](https://cmux.com/) (the default backend)
