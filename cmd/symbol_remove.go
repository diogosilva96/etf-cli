package cmd

import (
	"github.com/diogosilva96/etf-scraper/printer"
	"github.com/spf13/cobra"
)

func NewSymbolRemoveCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "remove",
		Aliases: []string{"r"},
		Short:   "Removes the specified ETF symbol from the tracked list.",
		Long: `Removes the specified ETF symbol to the tracked list.
		The remove ETF symbol will no longer be included in the fetched ETF data when the 'get' command is executed.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			symbol := args[0]
			err := removeSymbol(symbol)
			if err != nil {
				printer.PrintWarning(err.Error())
			}
		},
	}

	return &cmd
}

func removeSymbol(symbol string) error {
	err := c.RemoveSymbol(symbol)
	if err != nil {
		return err
	}
	err = c.Save(ConfigPath)
	if err != nil {
		return err
	}
	printer.Print("'%s' was successfully removed from the tracked list.", symbol)
	return nil
}
