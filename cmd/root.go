package cmd

import (
	"github.com/diogosilva96/etf-cli/internal/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "etf",
	Short: "etf - a simple CLI that generates real time ETF data reports & allows management of tracked ETFs.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute executes the command.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(func() {
		cobra.CheckErr(config.InitConfig())
	})
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
