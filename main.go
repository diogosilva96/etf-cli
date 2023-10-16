package main

import (
	"log"

	"github.com/diogosilva96/etf-scraper/app"
)

func main() {
	cfg, err := app.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	app.Run(cfg)
}
