package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/diogosilva96/etf-scraper/internal/cmd"
	"github.com/diogosilva96/etf-scraper/internal/config"
	"github.com/diogosilva96/etf-scraper/internal/printer"
)

const (
	ConfigPath = "config.json"
)

var Cfg config.Config

// Run runs the CLI application.
func Run(cfg *config.Config) {
	Cfg = *cfg
	if err := cmd.NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

// / TODO: fix this
func NewConfig() (*config.Config, error) {
	cfg, err := config.Parse(ConfigPath)
	if err != nil {
		printer.PrintWarning("The config in path '%s' could not be found or parsed.\nDetails: %s\nFalling back to default configuration.\n", ConfigPath, err.Error())
		cfg, err = config.NewConfig(config.WithSymbols("VWCE.DE", "VWCE.MI"))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error while falling back to default configuration.\nDetails: %s", err))
		}
		cfg.Save(ConfigPath)
	}
	return cfg, nil
}
