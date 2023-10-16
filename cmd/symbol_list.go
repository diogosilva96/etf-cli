package cmd

import (
	"github.com/diogosilva96/etf-scraper/printer"
	"github.com/spf13/cobra"
)

func NewSymbolListCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "Lists the ETF symbols in the tracked list.",
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			listSymbols()
		},
	}

	return &cmd
}

func listSymbols() {
	printer.Print("Tracked ETFs:\n")
	for _, s := range c.Symbols {
		printer.Print("- %s\n", s)
	}
}
