package cli

import (
	"fmt"
	"sort"

	"github.com/buemura/acm/internal/provider"
	"github.com/spf13/cobra"
)

var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Show supported providers",
	Run: func(cmd *cobra.Command, args []string) {
		names := provider.Names()
		sort.Strings(names)

		fmt.Println("Supported providers:")
		for _, name := range names {
			fmt.Printf("  - %s\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(providerCmd)
}
