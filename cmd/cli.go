package cmd

import (
	"github.com/diogosilva96/etf-cli/internal/config"
	"github.com/spf13/cobra"
)

// Cli represents a command line application.
type Cli struct {
	rootCmd cobra.Command
}

// NewCli creates a new cli application.
func NewCli() (*Cli, error) {
	err := config.InitConfig()
	if err != nil {
		return nil, err
	}
	return &Cli{*rootCmd}, nil
}

// Run executes the cli application.
func (c Cli) Run() error {
	return c.rootCmd.Execute()
}
