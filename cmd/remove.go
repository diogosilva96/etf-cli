package cmd

import (
	"github.com/diogosilva96/etf-cli/internal/config"
	"github.com/spf13/cobra"
	"strconv"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes an existing ETF from the configuration based on the specified index.",
	Long: `Removes an existing ETF from the configuration based on the specified index.

	When the 'report' command is used the data from the removed ETF will no longer be displayed.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idxStr := args[0]
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			cmd.PrintErrf("'%s' is not a valid number.\n", idxStr)
			return
		}
		etfs := config.ListEtfs()
		for i, etf := range etfs {
			if idx == i+1 {
				err := config.RemoveEtf(etf)
				if err != nil {
					cmd.PrintErr(err)
					return
				}
				cmd.Printf("'%s' successfully removed!", etf)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
