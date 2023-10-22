package cmd

import (
	"fmt"

	"github.com/diogosilva96/etf-cli/cmd/config"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes an existing ETF from the configuration.",
	Long: `Removes an existing ETF from the configuration.
	
	When the 'report' command is used the data from the removed ETF will no longer be displayed.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		etf := args[0]
		err := config.RemoveEtf(etf)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		cmd.Printf(fmt.Sprintf("etf '%s' successfully removed!", etf))
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
