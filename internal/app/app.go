package app

import (
	"sync"

	"github.com/diogosilva96/etf-scraper/internal/config"
	"github.com/diogosilva96/etf-scraper/internal/printer"
	"github.com/urfave/cli"
)

const (
	symbolArgName = "symbol"
	// The configuration's path.
	ConfigPath = "config.json"
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
			etfApp = &EtfApp{Config: initConfig()}
		}
	}
	return *etfApp
}

func initConfig() config.Config {
	c, err := config.Parse(ConfigPath)
	if err != nil {
		printer.PrintWarning("The config in path '%s' could not be found or parsed. Details: %s\nFalling back to default configuration.\n", ConfigPath, err.Error())
		c, err = config.NewConfig(config.WithSymbols("VWCE.DE", "VWCE.MI"))
		c.Save(ConfigPath)
	}
	return *c
}
g