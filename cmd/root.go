package cmd

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

var rootFlags struct {
	verbose bool
}

var rootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),
	Short: "Expose the UEC HPC laboratory Wiki through MCP",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		options := &tint.Options{AddSource: true, TimeFormat: time.TimeOnly}
		if rootFlags.verbose {
			options.Level = slog.LevelDebug
		}
		slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, options)))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&rootFlags.verbose, "verbose", "v", false, "Enable verbose output")
}
