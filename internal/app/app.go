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

// GetOrCreateEtfApp gets or creates a new etf application.
func GetOrCreateEtfApp() EtfApp {
	// TODO: This logic might make more sense to be in main.go - this also makes it so that we can remove the sync mechanism
	if etfApp == nil {
		lock.Lock()
		defer lock.Unlock()
		if etfApp == nil {
			app := cli.NewApp()
			app.Name = "ETF scraper CLI."
			app.Usage = "Let's you retrieve ETF information & manage tracked ETFs."
			app.Commands = []cli.Command{
				{
					Name:      "get",
					HelpName:  "get",
					Action:    getAction,
					ArgsUsage: "",
					Usage:     "Retrieves the information for all the ETFs in the tracked list.",
				},
				{
					Name:      "symbol",
					HelpName:  "symbol",
					ArgsUsage: "",
					Usage:     "Options for the tracked list symbols.",
					Subcommands: []cli.Command{
						{
							Name:      "list",
							HelpName:  "list",
							Aliases:   []string{"l"},
							Action:    listSymbolAction,
							ArgsUsage: "",
							Usage:     "Displays the list of tracked ETF symbols.",
						},
						{
							Name:        "add",
							HelpName:    "add",
							Aliases:     []string{"a"},
							Action:      addSymbolAction,
							ArgsUsage:   "",
							Usage:       "Adds the specified ETF symbol to the tracked list.",
							Description: "Add the ETF symbol to the tracked list. When the `get` command is run the ETF's data will be included.",
						},
						{
							Name:        "remove",
							HelpName:    "remove",
							Aliases:     []string{"r"},
							Action:      removeSymbolAction,
							ArgsUsage:   "",
							Usage:       "Removes the ETF symbol from the tracked list.",
							Description: "Removes the ETF from the tracked list. When the `get` command is run the removed ETF's data will no longer be included.",
						},
					},
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

	err := etfApp.get()
	if err != nil {
		return err
	}
	return nil
}

func listSymbolAction(c *cli.Context) error {
	if len(c.Args()) > 0 {
		return errors.New("No arguments were expected.")
	}
	etfApp.listSymbols()
	return nil
}

func addSymbolAction(c *cli.Context) error {
	if len(c.Args()) == 1 && !c.IsSet(symbolArgName) {
		symbol := c.Args()[0]
		err := etfApp.addSymbol(symbol)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Invalid number of arguments.")
}

func removeSymbolAction(c *cli.Context) error {
	if len(c.Args()) == 1 && !c.IsSet(symbolArgName) {
		symbol := c.Args()[0]
		err := etfApp.removeSymbol(symbol)
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("Invalid number of arguments.")
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

func (app EtfApp) addSymbol(symbol string) error {
	if app.Config.HasSymbol(symbol) {
		return errors.New(fmt.Sprintf("The symbol '%s' was not added to the tracked list, because it is already being tracked.", symbol))
	}
	app.Config.AddSymbol(symbol)
	err := app.Config.Save(configPath)
	if err != nil {
		return err
	}
	printer.Print("'%s' was successfully added to the tracked list.", symbol)
	return nil
}

func (app EtfApp) listSymbols() {
	printer.Print("Tracked ETFs:\n")
	for _, s := range app.Config.Symbols {
		printer.Print("- %s\n", s)
	}
}

func (app EtfApp) removeSymbol(symbol string) error {
	err := app.Config.RemoveSymbol(symbol)
	if err != nil {
		return err
	}
	err = app.Config.Save(configPath)
	if err != nil {
		return err
	}
	printer.Print("'%s' was successfully removed from the tracked list.", symbol)
	return nil
}

func (app EtfApp) get() error {

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
