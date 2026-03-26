package cli

import (
	"os"

	"github.com/buemura/acm/internal/ui"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{
	Use:     "acm",
	Short:   "Manage AI assistant cache files",
	Long:    "A CLI tool to manage AI assistant cache files, starting with Claude and extensible to other providers.",
	Version: version,
	Run:     rootRun,
}

func rootRun(cmd *cobra.Command, args []string) {
	ui.RunMenu()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
