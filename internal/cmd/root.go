package cmd

import (
	"github.com/diogosilva96/etf-scraper/internal/cmd/get"
	"github.com/diogosilva96/etf-scraper/internal/cmd/symbol"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "etf",
		Short: "etf - a simple CLI to retrieve ETF data reports & manage tracked ETFs.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(get.NewGetCmd(), symbol.NewSymbolRootCmd())

	return &cmd
}
