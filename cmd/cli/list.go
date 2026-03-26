package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/buemura/acm/internal/provider"
	"github.com/buemura/acm/internal/scanner"
	"github.com/buemura/acm/internal/ui"
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
	rootCmd.AddCommand(listCmd)
}
