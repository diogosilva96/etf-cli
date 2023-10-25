package cmd

import (
	"os"

	"github.com/diogosilva96/etf-cli/cmd/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "etf",
	Short: "etf - a simple CLI to real time ETF data reports & manage its configuration.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute executes the command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		rootCmd.PrintErr(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.InitConfig)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
