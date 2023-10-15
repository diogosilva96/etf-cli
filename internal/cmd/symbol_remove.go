package cmd

import (
	"github.com/diogosilva96/etf-scraper/internal/app"
	"github.com/diogosilva96/etf-scraper/internal/printer"
	"github.com/spf13/cobra"
)

var symbolRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"r"},
	Short:   "Removes the specified ETF symbol from the tracked list.",
	Long: `Removes the specified ETF symbol to the tracked list.
	The remove ETF symbol will no longer be included in the fetched ETF data when the 'get' command is executed.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		etfApp := app.GetOrCreateEtfApp() // TODO: rework this
		symbol := args[0]
		err := removeSymbol(etfApp, symbol)
		if err != nil {
			printer.PrintWarning(err.Error())
		}
	},
}

func removeSymbol(etfApp app.EtfApp, symbol string) error {
	err := etfApp.Config.RemoveSymbol(symbol)
	if err != nil {
		return err
	}
	err = etfApp.Config.Save(app.ConfigPath)
	if err != nil {
		return err
	}
	printer.Print("'%s' was successfully removed from the tracked list.", symbol)
	return nil
}

func init() {
	symbolRootCmd.AddCommand(symbolRemoveCmd)
}
