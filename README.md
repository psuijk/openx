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
openx add <name> [--path PATH]       # create a project config (defaults to cwd)
openx add-tab <name> <tab> [flags]   # add or update a tab
openx clone <source> <new> [--path]  # clone a project config
openx list                           # list all projects
openx show <name>                    # print a project's config
openx edit <name>                    # open config in $EDITOR
openx remove <name> [--yes]          # delete a project config
openx <name>                         # open the project workspace
openx <name> --dry-run               # print what would happen
openx <name> --join                  # new workspace in current window
openx <name> --new-window            # new workspace in a new window
openx <name> --backend <name>        # override backend
openx version                        # print version
```

Flags can go before or after the project name (e.g. `openx myproject --dry-run` or `openx --dry-run myproject`).

### Adding tabs

```bash
openx add-tab myproject claude --command claude    # add a tab named "claude"
openx add-tab myproject shell                      # add a tab with no command (plain shell)
openx add-tab myproject claude --command "claude --model opus"  # update existing tab's command
openx add-tab myproject logs --after claude         # insert after a specific tab
openx add-tab myproject logs --before shell          # insert before a specific tab
```

If a tab with the given name already exists, its command is updated in place without changing position.

### Cloning configs

```bash
openx clone myproject myproject2 --path /path/to/new/project  # clone with new path
openx clone myproject myproject2                               # clone, defaults path to cwd
```

Copies everything (tabs, pre/post open, backend, mode) from the source project.

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
