package cmd

import (
	"fmt"

	"github.com/diogosilva96/etf-scraper/cmd/config"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the ETFs in the configuration.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		etfs := config.ListEtfs()
		fmt.Printf("ETFs:\n")
		for _, e := range etfs {
			fmt.Printf("- %s\n", e)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
