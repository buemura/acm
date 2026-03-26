package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/buemura/agnt-cc/provider"
	"github.com/buemura/agnt-cc/scanner"
	"github.com/buemura/agnt-cc/ui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List cache files",
	Run: func(cmd *cobra.Command, args []string) {
		providerName, _ := cmd.Flags().GetString("provider")
		fileType, _ := cmd.Flags().GetString("type")
		age, _ := cmd.Flags().GetString("age")

		p, err := provider.Get(providerName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if missing := provider.MissingPathChecks(p); len(missing) > 0 {
			if len(missing) == len(p.Checks) {
				fmt.Fprintf(os.Stderr, "No %s cache found. Is %s installed?\n", p.Name, p.Name)
				return
			}
			fmt.Fprintf(
				os.Stderr,
				"Note: %s provider is missing expected %s in standard locations.\n",
				p.Name,
				strings.Join(missing, ", "),
			)
		}

		minAge, err := scanner.ParseAge(age)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		files, err := scanner.ScanFiles(p.CachePaths, scanner.FilterOpts{
			MinAge:   minAge,
			FileType: fileType,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		ui.PrintTable(files)
	},
}

func init() {
	listCmd.Flags().StringP("provider", "p", "claude", "Cache provider name")
	listCmd.Flags().StringP("type", "t", "", "Filter by file extension (e.g. json, log)")
	listCmd.Flags().StringP("age", "a", "", "Filter by minimum file age (e.g. 20m, 24h, 7d)")
	rootCmd.AddCommand(listCmd)
}
