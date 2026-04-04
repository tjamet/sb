package logging

import (
	"log/slog"
	"log/syslog"
	"os"
)

// EnableSyslog configures the default slog logger to write structured logs
// to the macOS system log (visible in Console.app under the "sb" subsystem).
// Falls back to stderr if the syslog connection fails.
func EnableSyslog() {
	w, err := syslog.New(syslog.LOG_DEBUG|syslog.LOG_USER, "sb")
	if err != nil {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
		return
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug})))
}
