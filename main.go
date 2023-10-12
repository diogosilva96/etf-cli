package main

import (
	"fmt"

	"github.com/diogosilva96/etf-scraper/internal/scraper"
)

const baseUrl = "https://finance.yahoo.com"

func main() {

	// TODO: make this a cli tool and allow removing/adding symbols
	symbols := []string{"VWCE.DE", "asd"}

	etfs := scraper.Scrape(symbols)

	fmt.Printf("%+v\n", etfs)
}
