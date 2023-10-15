package cmd

import (
	"errors"
	"fmt"

	"github.com/diogosilva96/etf-scraper/internal/app"
	"github.com/diogosilva96/etf-scraper/internal/printer"
	"github.com/spf13/cobra"
)

var symbolAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Adds the specified ETF symbol to the tracked list.",
	Long: `Adds the specified ETF symbol to the tracked list.
	The added ETF symbol will be included in the fetched ETF data when the 'get' command is executed.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		etfApp := app.GetOrCreateEtfApp() // TODO: rework this
		symbol := args[0]
		err := addSymbol(etfApp, symbol)
		if err != nil {
			printer.PrintWarning(err.Error())
		}
	},
}

func addSymbol(etfApp app.EtfApp, symbol string) error {
	if etfApp.Config.HasSymbol(symbol) {
		return errors.New(fmt.Sprintf("The symbol '%s' was not added to the tracked list, because it is already being tracked.", symbol))
	}
	etfApp.Config.AddSymbol(symbol)
	err := etfApp.Config.Save(app.ConfigPath)
	if err != nil {
		return err
	}
	printer.Print("'%s' was successfully added to the tracked list.", symbol)
	return nil
}

func init() {
	symbolRootCmd.AddCommand(symbolAddCmd)
}