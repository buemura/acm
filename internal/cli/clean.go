package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/buemura/agnt-cc/internal/provider"
	"github.com/buemura/agnt-cc/internal/scanner"
	"github.com/buemura/agnt-cc/internal/ui"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Delete cache files",
	Run: func(cmd *cobra.Command, args []string) {
		providerName, _ := cmd.Flags().GetString("provider")
		fileType, _ := cmd.Flags().GetString("type")
		age, _ := cmd.Flags().GetString("age")
		force, _ := cmd.Flags().GetBool("force")

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

		if len(files) == 0 {
			fmt.Println("No files found matching the criteria.")
			return
		}

		ui.PrintTable(files)

		total := scanner.TotalSize(files)
		fmt.Printf("\nSpace to be freed: %s\n", ui.FormatSize(total))

		if !force {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Delete these files? (y/N): ")
			input, _ := reader.ReadString('\n')
			if strings.ToLower(strings.TrimSpace(input)) != "y" {
				fmt.Println("Aborted.")
				return
			}
		}

		deleted := 0
		var freedBytes int64
		for _, f := range files {
			if err := os.Remove(f.Path); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to delete %s: %v\n", f.Path, err)
				continue
			}
			deleted++
			freedBytes += f.Size
		}

		fmt.Printf("Deleted %d file(s), freed %s\n", deleted, ui.FormatSize(freedBytes))
	},
}

func init() {
	cleanCmd.Flags().StringP("provider", "p", "claude", "Cache provider name")
	cleanCmd.Flags().StringP("type", "t", "", "Filter by file extension (e.g. json, log)")
	cleanCmd.Flags().StringP("age", "a", "", "Filter by minimum file age (e.g. 20m, 24h, 7d)")
	cleanCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	rootCmd.AddCommand(cleanCmd)
}
