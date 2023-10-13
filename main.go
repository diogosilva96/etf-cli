package main

import (
	"fmt"

	"github.com/diogosilva96/etf-scraper/internal/config"
	"github.com/diogosilva96/etf-scraper/internal/printer"
	"github.com/diogosilva96/etf-scraper/internal/scraper"
)

const (
	configPath = "config.json"
)

var cfg config.Config
var defaultCfg config.Config = config.Config{Symbols: []string{"VWCE.DE"}}

func main() {

	// TODO: make this a cli tool and allow removing/adding symbols

	cfg, err := config.Parse(configPath)

	if err != nil {
		printer.PrintWarning("The config in path '%s' could not be found or parsed. Details: %s\nFalling back to default configuration.\n", configPath, err.Error())
		cfg = &defaultCfg
		cfg.Save(configPath)
	}

	fmt.Printf("%+v\n", cfg)

	etfs := scrape(cfg)

	fmt.Printf("\n%+v\n", etfs)
}

func scrape(cfg *config.Config) []scraper.Etf {
	printer.Print("Scraping data...\n")
	etfs := make([]scraper.Etf, 0, len(cfg.Symbols))
	for _, s := range cfg.Symbols {
		etf, err := scraper.Scrape(s)
		if err != nil {
			printer.PrintWarning("[%s] Something went wrong when scraping the data. Error details: %s\n", s, err)
			continue
		}
		etfs = append(etfs, *etf)
		printer.PrintInfo("[%s] Success!\n", s)
	}
	fmt.Print("Scraping complete.\n")
	return etfs
}
