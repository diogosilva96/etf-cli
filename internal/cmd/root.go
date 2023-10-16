package cmd

import (
	"github.com/diogosilva96/etf-scraper/internal/cmd/get"
	"github.com/diogosilva96/etf-scraper/internal/cmd/symbol"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "etf",
		Short: "etf - a simple CLI to retrieve ETF data & manage ETF tracked list.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cmd.AddCommand(get.NewGetCmd(), symbol.NewSymbolRootCmd())

	return &cmd
}
