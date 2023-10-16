package cmd

import (
	"github.com/spf13/cobra"
)

func NewSymbolRootCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "symbol",
		Short: "Options for managing the ETF symbols in the tracked list.",
		Long:  `Options for adding, listing or removing symbols from the tracked list.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cmd.AddCommand(NewSymbolAddCmd(), NewSymbolRemoveCmd(), NewSymbolListCmd())

	return &cmd
}
