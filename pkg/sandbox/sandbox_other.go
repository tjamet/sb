//go:build !darwin

package sandbox

import (
	"errors"
	"runtime"
)

// Run is unsupported off macOS: the sandbox is built on sandbox-exec.
func Run(_ *Config, _ []string) error {
	return errors.New("sandbox-exec is macOS only, not available on " + runtime.GOOS)
}
