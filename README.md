# openx

A CLI tool for opening project-specific terminal workspaces. Run `openx <project-name>` to spin up your full dev environment — terminal tabs, editor, background tasks — all from a single command.

## Install

```bash
go build -o openx ./cmd/openx
```

## Usage

```
openx add <name> [--path PATH]   # create a project config (defaults to cwd)
openx list                       # list all projects
openx show <name>                # print a project's config
openx edit <name>                # open config in $EDITOR
openx remove <name> [--yes]      # delete a project config
openx <name>                     # open the project workspace
openx version                    # print version
```

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

## Status

Early development. Config management works (add, list, show, edit, remove, validation). Backend execution (actually opening workspaces) is next.
