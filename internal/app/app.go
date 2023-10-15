package app

import (
	"errors"
	"fmt"
	"sync"

	"github.com/diogosilva96/etf-scraper/internal/config"
	"github.com/diogosilva96/etf-scraper/internal/printer"
	"github.com/diogosilva96/etf-scraper/internal/scraper"
	"github.com/urfave/cli"
)

const (
	symbolArgName = "symbol"
	configPath    = "config.json"
)

// Represents an ETF application.
type EtfApp struct {
	*cli.App
	Config config.Config
}

var lock = &sync.Mutex{}

var etfApp *EtfApp

// NewEtfApp creates a new etf application.
func NewEtfApp() EtfApp {
	if etfApp == nil {
		lock.Lock()
		defer lock.Unlock()
		if etfApp == nil {
			app := cli.NewApp()
			app.Name = "ETF scraper CLI."
			app.Usage = "Let's you get, track & untrack ETFs."
			app.Commands = []cli.Command{
				{
					Name:        "track",
					HelpName:    "track",
					Action:      trackAction,
					ArgsUsage:   "",
					Usage:       "Adds the ETF to the tracked list based on its symbol.",
					Description: "Starts tracking the ETF. When the `get` command is run this ETF's data will also be retrieved.",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     symbolArgName,
							Usage:    "track ETF by symbol. ",
							Required: false,
						},
					},
				},
				{
					Name:        "get",
					HelpName:    "get",
					Action:      getAction,
					ArgsUsage:   "",
					Usage:       "Retrieves the information for all the ETFs in the tracked list.",
					Description: "Retrieves the information for all the ETFs in the tracked list.",
				},
			}

			etfApp = &EtfApp{App: app, Config: initConfig()}
		}
	}
	return *etfApp
}

func getAction(c *cli.Context) error {
	if len(c.Args()) > 0 {
		return errors.New("No arguments were expected.")
	}

	etfs := scrape(*etfApp)

	fmt.Printf("\n%+v\n", etfs)
	return nil
}

func scrape(app EtfApp) []scraper.Etf {
	type result struct {
		symbol string
		etf    *scraper.Etf
		err    error
	}

	ch := make(chan result, len(app.Config.Symbols))
	wg := sync.WaitGroup{}
	printer.Print("Scraping data...\n")
	for _, s := range app.Config.Symbols {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			etf, err := scraper.Scrape(symbol)
			r := result{symbol: symbol, etf: etf, err: err}
			ch <- r
		}(s)
	}
	etfs := make([]scraper.Etf, 0, len(app.Config.Symbols))
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

func trackAction(c *cli.Context) error {
	if len(c.Args()) == 1 && !c.IsSet(symbolArgName) {
		symbol := c.Args()[0]
		err := track(*etfApp, symbol)
		if err != nil {
			return err
		}
		return nil
	}

	if c.IsSet(symbolArgName) {
		symbol := c.String(symbolArgName)
		err := track(*etfApp, symbol)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Invalid number of arguments.")
}

func track(app EtfApp, symbol string) error {
	if app.Config.HasSymbol(symbol) {
		return errors.New(fmt.Sprintf("The symbol '%s' was not added to the tracked list, because it is already being tracked.", symbol))
	}
	app.Config.AddSymbol(symbol)
	err := app.Config.Save(configPath)
	if err != nil {
		return err
	}
	return nil
}

func initConfig() config.Config {
	c, err := config.Parse(configPath)
	if err != nil {
		printer.PrintWarning("The config in path '%s' could not be found or parsed. Details: %s\nFalling back to default configuration.\n", configPath, err.Error())
		c, err = config.NewConfig(config.WithSymbols("VWCE.DE", "VWCE.MI"))
		c.Save(configPath)
	}
	return *c
}
