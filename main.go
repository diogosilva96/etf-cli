package main

import (
	"log"
	"os"

	"github.com/diogosilva96/etf-scraper/internal/app"
)

const (
	configPath = "config.json"
)

func main() {
	app := app.GetOrCreateEtfApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
