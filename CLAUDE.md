# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Critical rule

**Never implement code.** This is a learning project. Guide, explain, and suggest — but the user writes all the code themselves. Exception: mechanical tasks with no learning value (renaming, moving files, find-and-replace) — just do those directly.

## What is openx

A CLI tool for opening project-specific terminal workspaces. Run `openx <project-name>` and your full dev environment spins up: cmux workspace with tabs, VS Code, etc. Configs are TOML files stored per-project under `$XDG_CONFIG_HOME/openx/projects/`.

## Build and run

```bash
go build -o openx ./cmd/openx     # build the binary
go run ./cmd/openx                 # run without building
go test ./...                      # run all tests
go test ./internal/config/         # run tests for one package
go vet ./...                       # static analysis
```

No external dependencies yet — stdlib `flag` for CLI, no third-party libraries.

## Architecture

- **`cmd/openx/main.go`** — tiny entry point, delegates to `command.Dispatch`
- **`internal/command/`** — subcommand handlers (add, list, show, mod, remove, run) and `dispatch.go` which routes `os.Args` to the right handler. Handlers should be dumb glue (~40 lines max); extract logic to other packages
- **`internal/config/`** — (planned) TOML config schema, XDG path resolution, validation
- **`internal/backend/`** — (planned) pluggable `Backend` interface with cmux as first implementation
- **`internal/plan/`** — (planned) shared `Plan` type with pre_open/post_open steps + opaque backend steps
- **`internal/shell/`** — (planned) `os/exec` wrapper for running shell commands
- **`internal/editor/`** — (planned) `$EDITOR` shell-out for config editing

## Design decisions (from docs/DESIGN.md)

- **Subcommands, not flags** — `openx add`, `openx list`, bare `openx <name>` routes to run
- **TOML, one file per project** — stored at `~/.config/openx/projects/<name>.toml`
- **Pluggable backends** — `Backend` interface; cmux first, tmux later. New backend = new subdir under `internal/backend/`
- **stdlib `flag` only** — each subcommand owns its own `flag.FlagSet`
- **pre_open/post_open** — first-class side actions (e.g. `code .`, `git fetch`), not special-cased features
- **Mode resolution** — `default_mode` in config, overridable by `--join`/`--new-window` flags; inside-cmux detection via `CMUX_WINDOW_ID` env var

## Current state

Early scaffolding (Milestone 0 from the build plan). Dispatch works, all subcommand handlers are stubs returning nil. The build plan in `docs/DESIGN.md` has 5 milestones — currently working through M0.

## Adding new things

- New subcommand: add handler file in `internal/command/`, add case in `dispatch.go`
- New backend: new subdir under `internal/backend/`, register in `registry.go`
- New config field: add to struct in `config/`, validate in `validate.go`
