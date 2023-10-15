package cmd

import (
	"github.com/diogosilva96/etf-scraper/internal/app"
	"github.com/diogosilva96/etf-scraper/internal/printer"
	"github.com/spf13/cobra"
)

var symbolListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Lists the ETF symbols in the tracked list.",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		etfApp := app.GetOrCreateEtfApp() // TODO: rework this
		listSymbols(etfApp)
	},
}

func listSymbols(etfApp app.EtfApp) {
	printer.Print("Tracked ETFs:\n")
	for _, s := range etfApp.Config.Symbols {
		printer.Print("- %s\n", s)
	}
}

func init() {
	symbolRootCmd.AddCommand(symbolListCmd)
}
