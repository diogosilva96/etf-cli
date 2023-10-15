package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/diogosilva96/etf-scraper/internal/cmd"
	"github.com/diogosilva96/etf-scraper/internal/config"
	"github.com/diogosilva96/etf-scraper/internal/printer"
)

var cfg *config.Config

func main() {

	cfg, err := initConfig()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Run(cfg)
}

func initConfig() (*config.Config, error) {
	cfg, err := config.Parse(cmd.ConfigPath)
	if err != nil {
		printer.PrintWarning("The config in path '%s' could not be found or parsed.\nDetails: %s\nFalling back to default configuration.\n", cmd.ConfigPath, err.Error())
		cfg, err = config.NewConfig(config.WithSymbols("VWCE.DE", "VWCE.MI"))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error while falling back to default configuration.\nDetails: %s", err))
		}
		cfg.Save(cmd.ConfigPath)
	}
	return cfg, nil
}
