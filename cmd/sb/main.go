package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tjamet/sb/pkg/logging"
	"github.com/tjamet/sb/pkg/sandbox"
)

func main() {
	var verbose bool

	rootCmd := &cobra.Command{
		Use:   "sb [flags] <command> [args...]",
		Short: "Run a command in a macOS sandbox",
		Long: `sb runs a command inside a macOS sandbox-exec profile with restricted
filesystem access. Permissions are configured via YAML files stored
in ~/.config/sb/<command>.yaml.`,
		Args:              cobra.MinimumNArgs(1),
		DisableFlagParsing: false,
		SilenceUsage:      true,
		TraverseChildren:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				logging.EnableSyslog()
			}
			cfg, err := sandbox.FromCommand(args[0])
			if err != nil {
				return fmt.Errorf("loading configuration: %w", err)
			}
			return sandbox.Run(cfg, args[1:])
		},
	}
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable debug logging to macOS system log (Console.app)")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
