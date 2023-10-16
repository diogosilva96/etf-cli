package cmd

import (
	"errors"
	"fmt"

	"github.com/diogosilva96/etf-scraper/printer"
	"github.com/spf13/cobra"
)

func NewSymbolAddCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "Adds the specified ETF symbol to the tracked list.",
		Long: `Adds the specified ETF symbol to the tracked list.
		The added ETF symbol will be included in the fetched ETF data when the 'get' command is executed.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\n%+v\n", args)
			symbol := args[0]
			err := addSymbol(symbol)
			if err != nil {
				printer.PrintWarning(err.Error())
			}
		},
	}

	return &cmd
}

func addSymbol(symbol string) error {
	if c.HasSymbol(symbol) {
		return errors.New(fmt.Sprintf("The symbol '%s' was not added to the tracked list, because it is already being tracked.", symbol))
	}
	c.AddSymbol(symbol)
	err := c.Save(ConfigPath)
	if err != nil {
		return err
	}
	printer.Print("'%s' was successfully added to the tracked list.", symbol)
	return nil
}