package cmd

import (
	"fmt"
	"os"

	"github.com/diogosilva96/etf-scraper/internal/config"
	"github.com/spf13/cobra"
)

const (
	ConfigPath = "config.json"
)

// Represents the configuration for the cmd module.
var c *config.Config

var rootCmd = &cobra.Command{
	Use:   "etf",
	Short: "etf - a simple CLI to retrieve ETF data & manage ETF tracked list.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Run runs the CLI application.
func Run(cfg *config.Config) {
	c = cfg
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
