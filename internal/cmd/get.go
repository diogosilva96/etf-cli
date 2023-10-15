package cmd

import (
	"fmt"
	"sync"

	"github.com/diogosilva96/etf-scraper/internal/app"
	"github.com/diogosilva96/etf-scraper/internal/printer"
	"github.com/diogosilva96/etf-scraper/internal/scraper"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve ETF data based on the tracked ETF symbols.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		etfApp := app.GetOrCreateEtfApp() // TODO: rework this
		etfs := scrape(etfApp)

		fmt.Printf("\n%+v\n", etfs)
	},
}

func scrape(app app.EtfApp) []scraper.Etf {
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

func init() {
	rootCmd.AddCommand(getCmd)
}
