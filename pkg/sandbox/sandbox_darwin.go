package sandbox

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func darwinSandboxPermissions(perm Permission) []string {
	var permissions []string
	if perm&ReadMeta != 0 {
		permissions = append(permissions, "file-read-metadata")
	}
	if perm&ReadData != 0 {
		permissions = append(permissions, "file-read-data")
	}
	if perm&WriteMeta != 0 {
		permissions = append(permissions, "file-write-metadata")
	}
	if perm&WriteData != 0 {
		permissions = append(permissions, "file-write-data")
	}
	return permissions
}

func Run(cfg *Config, args []string) error {
	slog.Debug("running sandbox", "command", cfg.Command, "args", args, "cwd", cfg.cwd)

	sandboxInstructions := []string{
		"(version 1)",
	}

	sandboxPath := filepath.Join(cfg.configDir, cfg.commandName+".sb")
	slog.Debug("checking for sandbox profile", "path", sandboxPath)
	_, err := os.Stat(sandboxPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("checking sandbox profile %s: %w", sandboxPath, err)
		}
		slog.Debug("no custom sandbox profile found", "path", sandboxPath)
	} else {
		slog.Debug("importing custom sandbox profile", "path", sandboxPath)
		sandboxInstructions = append(sandboxInstructions, fmt.Sprintf(`(import "%s")`, sandboxPath))
	}

	if len(cfg.FoldersInParents) > 0 {
		perms := darwinSandboxPermissions(cfg.IncludeParents)
		slog.Debug("resolved parent permissions", "permissions", perms)

		parent := cfg.cwd
		for {
			parent = filepath.Dir(parent)
			if len(perms) > 0 {
				sandboxInstructions = append(sandboxInstructions, fmt.Sprintf(`(allow %s (literal "%s"))`, strings.Join(perms, " "), parent))
			}
			for _, folderPerm := range cfg.FoldersInParents {
				folderPath := filepath.Join(parent, folderPerm.Path)
				folderPerms := darwinSandboxPermissions(folderPerm.Permission)
				if len(folderPerms) > 0 {
					sandboxInstructions = append(sandboxInstructions, fmt.Sprintf(`(allow %s (subpath "%s"))`, strings.Join(folderPerms, " "), folderPath))
				}
			}
			if parent == "/" || parent == "" || parent == cfg.home {
				break
			}
		}
	}

	slog.Debug("generated sandbox profile", "instructions", strings.Join(sandboxInstructions, "\n"))

	cmd := append([]string{cfg.Command}, args...)
	if len(sandboxInstructions) > 1 {
		sandboxExec := []string{"sandbox-exec", "-p", strings.Join(sandboxInstructions, "\n")}
		for _, env := range cfg.ExportEnvironment {
			sandboxExec = append(sandboxExec, "-D"+env+"="+os.Getenv(env))
		}
		cmd = append(sandboxExec, cmd...)
	}

	pgm := cmd[0]
	if !filepath.IsAbs(pgm) {
		targetProgram, err := exec.LookPath(pgm)
		if err != nil {
			return fmt.Errorf("program not found %s: %w", pgm, err)
		}
		pgm = targetProgram
	}

	slog.Debug("executing command", "program", pgm, "argv", cmd)
	return syscall.Exec(pgm, cmd, os.Environ())
}
