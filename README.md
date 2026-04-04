# sb

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A macOS sandbox wrapper that runs commands with restricted filesystem access using `sandbox-exec`.

`sb` loads a YAML configuration for the given command, generates a sandbox profile that limits file access to specific directories, and executes the command within that sandbox.

## Installation

```sh
go install github.com/tjamet/sb/cmd/sb@latest
```

## Usage

```sh
sb <command> [args...]
```

### Examples

Run `git status` in a sandbox:

```sh
sb git status

# alternatively, you can alias to git
alias git="sb git"
```

Run with debug logging enabled:

```sh
env SB_LOG_LEVEL=info sb git status
```

## Configuration

Configuration files are stored at `~/.config/sb/<command>.yaml`. When you run `sb git ...`, it looks for `~/.config/sb/git.yaml`.

Checkout the [examples](./examples) for real-world examples

### Config file format

```yaml
# The command to execute (overrides the command name from the CLI)
command: /usr/bin/git

# Permissions granted to all parent directories from the working directory up to $HOME
include_parents:
  - read-meta
  - read-data

# Additional folder permissions applied at each parent level
folders_in_parents:
  - path: .git
    permissions:
      - read-meta
      - read-data
      - write-meta
      - write-data
  - path: node_modules
    permissions:
      - read-meta
      - read-data

# Environment variables to export into the sandbox
export_env:
  - HOME
  - PATH
  - PWD
```

### Permission types

| Permission | sandbox-exec equivalent | Description |
|------------|------------------------|-------------|
| `read-meta` | `file-read-metadata` | Read file metadata (stat, listing) |
| `read-data` | `file-read-data` | Read file contents |
| `write-meta` | `file-write-metadata` | Modify file metadata (rename, chmod) |
| `write-data` | `file-write-data` | Write file contents |

### How permissions are applied

When `folders_in_parents` is configured, `sb` walks up from the current working directory to `$HOME` (or `/`). At each parent directory:

1. The `include_parents` permissions are granted on the parent directory itself.
2. Each entry in `folders_in_parents` is joined with the parent path, and its permissions are granted on that subdirectory.

This allows patterns like granting read access to project directories while granting write access only to specific subdirectories (e.g., `.git`).

### Custom sandbox profiles

You can provide a custom `sandbox-exec` profile at `~/.config/sb/<command>.sb`. If this file exists, it is imported into the generated profile with `(import ...)`, allowing you to add rules that `sb` does not generate (e.g., network access, process execution).

Example `~/.config/sb/git.sb`:

```scheme
(allow network-outbound (remote tcp))
(allow process-exec)
```

These are good source of information to write the profiles:
- https://gibbok.github.io/myvar/macos/macoss-command-line-sandboxing-tool-sandbox-exec/
- https://gist.github.com/cedws/b7dd60878a208be09ae7d766e6317061

### Debugging the profiles

You can debug the profiles to understand what are the blocked operations using the macos console:

```sh
log stream --style=compact --predicate 'sender=="Sandbox"'
```

## Debugging sb

Logs are always enabled with the `SB_LOG_LEVEL` environment variable. View them in Console.app by filtering for the process `sb`:

```sh
export SB_LOG_LEVEL=debug # | info | warn | error
sb git status
```

Or from the terminal:

```sh
log stream --predicate 'process == "sb"' --level debug
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
