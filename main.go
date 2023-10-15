package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/diogosilva96/etf-scraper/internal/config"
	"github.com/diogosilva96/etf-scraper/internal/printer"
	"github.com/diogosilva96/etf-scraper/internal/scraper"
	"github.com/urfave/cli"
)

const (
	configPath    = "config.json"
	symbolArgName = "symbol"
)

var cfg config.Config

func main() {

	fmt.Printf("\n%+v\n", cfg)
	// TODO: make this a cli tool and allow removing/adding symbols
	app := cli.NewApp()
	app.Name = "ETF scraper CLI."
	app.Usage = "Let's you fetch, track & untrack ETFs."
	app.Commands = []cli.Command{
		{
			Name:        "track",
			HelpName:    "track",
			Action:      trackAction,
			ArgsUsage:   "",
			Usage:       "Adds the ETF to the tracked list based on its symbol.",
			Description: "Starts tracking the ETF. When the `fetch` command is run this ETF's data will also be fetched.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     symbolArgName,
					Usage:    "track ETF by symbol. ",
					Required: false,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: add way to validate if an etf being added is valid (scrape to see if it can be found by the symbol)
	etfs := scrape(&cfg)

	fmt.Printf("\n%+v\n", etfs)
}

func trackAction(c *cli.Context) error {
	if len(c.Args()) == 1 && !c.IsSet(symbolArgName) {
		symbol := c.Args()[0]
		err := track(symbol)
		if err != nil {
			return err
		}
		return nil
	}

	if c.IsSet(symbolArgName) {
		symbol := c.String(symbolArgName)
		err := track(symbol)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Invalid number of arguments.")
}

func track(symbol string) error {
	if cfg.Contains(symbol) {
		return errors.New(fmt.Sprintf("The symbol '%s' was not added to the tracked list, because it is already being tracked.", symbol))
	}
	cfg.Add(symbol)
	cfg.Save(configPath)
	return nil
}

func scrape(cfg *config.Config) []scraper.Etf {
	type result struct {
		symbol string
		etf    *scraper.Etf
		err    error
	}

	ch := make(chan result, len(cfg.Symbols))
	wg := sync.WaitGroup{}
	printer.Print("Scraping data...\n")
	for _, s := range cfg.Symbols {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			etf, err := scraper.Scrape(symbol)
			r := result{symbol: symbol, etf: etf, err: err}
			ch <- r
		}(s)
	}
	etfs := make([]scraper.Etf, 0, len(cfg.Symbols))
	wg.Wait()
	close(ch)

	for r := range ch {
		if r.err != nil {
			printer.PrintWarning("[%s] Something went wrong while scraping the data. Error details: %s\n", r.symbol, r.err)
			continue
		}
		etfs = append(etfs, *r.etf)
		printer.PrintInfo("[%s] Success!\n", r.symbol)
	}
	printer.Print("Scraping complete.\n")
	return etfs
}

func initConfig(filePath string) config.Config {
	c, err := config.Parse(filePath)
	if err != nil {
		printer.PrintWarning("The config in path '%s' could not be found or parsed. Details: %s\nFalling back to default configuration.\n", filePath, err.Error())
		c, err = config.NewConfig(config.WithSymbols("VWCE.DE", "VWCE.MI"))
		c.Save(filePath)
	}
	return *c
}

func init() {
	cfg = initConfig(configPath)
}
