# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`sb` is a macOS sandbox wrapper that runs commands with restricted filesystem access using `sandbox-exec`. It loads YAML config for a given command (`~/.config/sb/<command>.yaml`), generates a sandbox profile limiting file access, and executes the command within that sandbox.

macOS-only (`sandbox-exec`, syslog). Platform-specific files use `_darwin.go` suffix.

## Build & Test

```sh
go build ./...          # build all
go build ./cmd/sb       # build the binary
go test ./...           # run all tests
go test ./pkg/sandbox   # run tests for a specific package
go test -run TestConfigMarshal ./pkg/sandbox  # run a single test
```

## Architecture

- `cmd/sb/main.go` — CLI entrypoint using cobra. Parses flags, loads config via `sandbox.FromCommand()`, runs via `sandbox.Run()`.
- `pkg/sandbox/config.go` — Config loading from YAML, `Permission` bitmask type with YAML marshal/unmarshal.
- `pkg/sandbox/sandbox_darwin.go` — Sandbox profile generation (sandbox-exec S-expression format) and command execution via `syscall.Exec`. Walks parent directories from cwd to `$HOME` applying folder permissions.
- `pkg/logging/logging_darwin.go` — Optional syslog-based debug logging (Console.app).

The sandbox profile is built as a list of S-expression instructions, optionally importing a custom `.sb` file from the config directory.

## Conventions

- Build with plain `go build` (no special flags or Makefile).
- Use `github.com/stretchr/testify/assert` and `require` for test assertions.
- Prefer test helpers over table-driven tests with nested loops. Keep test output readable.
- Mock dependencies with hand-written interface implementations (struct with `FooFunc` field and `FooCalls` counter), not mocking libraries.
- Use `github.com/spf13/cobra` for CLI commands.
- YAML library is `go.yaml.in/yaml/v3` (not `gopkg.in/yaml.v3` directly).
