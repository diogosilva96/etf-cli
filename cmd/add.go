package cmd

import (
	"fmt"

	"github.com/diogosilva96/etf-scraper/cmd/config"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new ETF to the configuration.",
	Long: `Adds a new ETF to the configuration.
	
	When the 'report' command is used all the data for the ETFs in the configuration will be displayed.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		etf := args[0]
		err := config.AddEtf(etf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf(fmt.Sprintf("etf '%s' successfully added!", etf))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
