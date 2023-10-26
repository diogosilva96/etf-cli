package cmd

import (
	"github.com/diogosilva96/etf-cli/internal/config"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the ETFs in the configuration.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		etfs := config.ListEtfs()
		if len(etfs) == 0 {
			cmd.Println("There are no etfs in the configuration.")
			return
		}
		cmd.Printf("ETFs:\n")
		for i, e := range etfs {
			cmd.Printf(" %v. %s\n", i+1, e)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
