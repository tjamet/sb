package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/tjamet/sb/pkg/sandbox"
)

func main() {

	logWriter := io.Discard
	logLevel := slog.LevelError

	switch os.Getenv("SB_LOG_LEVEL") {
	case "":
	case "debug":
		logWriter = os.Stderr
		logLevel = slog.LevelDebug
	case "info":
		logWriter = os.Stderr
		logLevel = slog.LevelInfo
	case "warn":
		logWriter = os.Stderr
		logLevel = slog.LevelWarn
	case "error":
		logWriter = os.Stderr
		logLevel = slog.LevelError
	default:
		fmt.Fprintf(os.Stderr, "invalid log level: %s (expected: debug, info, warn, error)\n", os.Getenv("SB_LOG_LEVEL"))
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(logWriter, &slog.HandlerOptions{Level: logLevel})))

	args := os.Args

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <command> [args...]\n", args[0])
		os.Exit(1)
	}

	cfg, err := sandbox.FromCommand(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	err = sandbox.Run(cfg, args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
