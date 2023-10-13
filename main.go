package main

import (
	"fmt"

	"github.com/diogosilva96/etf-scraper/internal/config"
	"github.com/diogosilva96/etf-scraper/internal/printer"
	"github.com/diogosilva96/etf-scraper/internal/scraper"
)

const baseUrl = "https://finance.yahoo.com"

func main() {

	// TODO: make this a cli tool and allow removing/adding symbols
	symbols := []string{"VWCE.DE", "asd"}

	config, err := config.Parse("config.json")

	if err != nil {
		printer.PrintError(err.Error())
		return
	}

	fmt.Printf("%+v\n", config)

	etfs := scraper.Scrape(symbols)

	fmt.Printf("%+v\n", etfs)
}
