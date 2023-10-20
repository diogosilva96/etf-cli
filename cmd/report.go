package cmd

import (
	"fmt"
	"sync"

	"github.com/diogosilva96/etf-scraper/cmd/config"
	"github.com/diogosilva96/etf-scraper/cmd/scraper"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Provides a report containing the up to date information of the ETFs in the tracked list.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		generateReport()
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}

func generateReport() {
	type result struct {
		symbol string
		etf    *scraper.Etf
		err    error
	}
	etfs := config.ListEtfs()
	ch := make(chan result, len(etfs))
	wg := sync.WaitGroup{}
	fmt.Printf("Scraping data...\n")
	for _, s := range etfs {
		wg.Add(1)
		go func(symbol string) {
			defer wg.Done()
			etf, err := scraper.Scrape(symbol)
			r := result{symbol: symbol, etf: etf, err: err}
			ch <- r
		}(s)
	}
	etfData := make([]scraper.Etf, 0, len(etfs))
	wg.Wait()
	close(ch)

	for r := range ch {
		if r.err != nil {
			fmt.Printf("[%s] Something went wrong while scraping the data. Error details: %s\n", r.symbol, r.err)
			continue
		}
		etfData = append(etfData, *r.etf)
		fmt.Printf("[%s] Success!\n", r.symbol)
	}
	fmt.Printf("Scraping complete.\n")

	fmt.Printf("%+v\n", etfData)
}
