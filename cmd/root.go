package cmd

import (
	"fmt"
	"os"

	"github.com/diogosilva96/etf-scraper/config"
	"github.com/spf13/cobra"
)

const (
	ConfigPath = "config.json"
)

// Represents the configuration for the cmd module.
var c *config.Config

func NewRootCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "etf",
		Short: "etf - a simple CLI to retrieve ETF data & manage ETF tracked list.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cmd.AddCommand(NewGetCmd(), NewSymbolRootCmd())

	return &cmd
}

// Run runs the CLI application.
func Run(cfg *config.Config) {
	c = cfg
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
