package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "0.1.1"

var rootCmd = &cobra.Command{
	Use:     "acm",
	Short:   "Manage AI assistant cache files",
	Long:    "A CLI tool to manage AI assistant cache files.",
	Version: version,
	Run:     rootRun,
}

func rootRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func init() {
	rootCmd.PersistentFlags().StringP("provider", "p", "claude", "Cache provider name")
	rootCmd.PersistentFlags().StringP("type", "t", "", "Filter by file extension (e.g. json, log)")
	rootCmd.PersistentFlags().StringP("age", "a", "", "Filter by minimum file age (e.g. 20m, 24h, 7d)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
