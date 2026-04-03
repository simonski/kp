# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`kp` is a terminal key/pair manager written in Go. It stores encrypted key/value entries in a JSON file (`~/.kpfile`) using RSA key-based encryption (`~/.ssh/kp.id_rsa`). It supports a TUI mode via termui (`kp -g` or `KP_GUI=1`).

## Build & Development Commands

```bash
make setup              # Install dependencies (staticcheck, bn)
make build              # Lint, bump version, format, build native binary
make test               # Run tests: go test
make install            # go install
make release            # Cross-compile for darwin/linux × amd64/arm64, create GitHub release via gh
make release-formula    # Generate Homebrew formula from template, push to simonski/homebrew-tap via gh API
go test -run TestMain   # Run a single test
```

Release flow: `make release` then `make release-formula`. The formula template lives in `homebrew/kp.rb.tmpl`.

## Architecture

All source is in a single `main` package (no subdirectories):

- **main.go** — CLI entry point and all command handlers (`DoGet`, `DoPut`, `DoList`, `DoDelete`, `DoRename`, etc.). Command dispatch is a chain of `is*()` checks in `main()`. Password generation uses `crypto/rand` with a controlled character set (`AllowedChars`).
- **objects.go** — Data model: `KPDB` (load/save/encrypt/decrypt), `DB` (versioned JSON structure with `Entries` and `History` maps), `DBEntry` (key, value, metadata, tags). Encryption delegates to `github.com/simonski/goutils/crypto`. The `Load` method handles schema migration from pre-versioned (plain map) to versioned DB format.
- **constants.go** — Environment variable names (`KP_FILE`, `KP_KEY`), default paths, usage text, SSH keygen command template.
- **gui.go** — Full-screen TUI using `gizak/termui/v3`. Q quits, Enter copies value to clipboard, E enters edit mode.

## Key Design Details

- Values are encrypted at rest in the JSON file; they are decrypted on read via `GetDecrypted()` and re-encrypted on write via `Put()`.
- `Put()` in objects.go has inverted exists/!exists logic for `Created` timestamps — new entries get `time.Now()`, existing entries preserve their original `Created`.
- `DoPut` supports three value sources: `-value` flag, `-random <size>` for generated passwords, or interactive stdin password prompt (no echo).
- The `Buildnumber` file is embedded via `//go:embed` and used as the version string.

## Key Dependencies

- `github.com/simonski/cli` — CLI argument parsing
- `github.com/simonski/goutils` — File utilities, string padding, crypto wrappers
- `github.com/atotto/clipboard` — Clipboard access for `get` command
- `github.com/gizak/termui/v3` — Terminal UI
- `github.com/pkg/browser` — Opens URLs for `open` command

## Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `KP_FILE` | `~/.kpfile` | Path to the encrypted key/pair database |
| `KP_KEY` | `~/.ssh/kp.id_rsa` | Path to RSA private key for encryption |
| `KP_GUI` | `0` | Set to `1` to launch TUI mode |

## Task Tracking

Tickets are managed with `tk` (install via `brew install simonski/tap/ticket`). Run `tk list` to see open work. The `.ticket/` directory holds the local ticket database.

## CI

GitHub Actions workflows in `.github/workflows/`: `go_again.yml` (compile and test). Releases are done locally via `make publish`.

- use red/green testing
- do not push if any tests are failing