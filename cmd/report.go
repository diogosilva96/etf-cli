package cmd

import (
	"fmt"
	"sync"

	"github.com/diogosilva96/etf-cli/cmd/config"
	"github.com/diogosilva96/etf-cli/cmd/report"
	"github.com/diogosilva96/etf-cli/cmd/scraper"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Provides a report containing the up to date information of the ETFs in the configuration.",
	Long: `Provides a report containing the up to date information of the ETFs in the configuration.
	
	A report will be generated for each ETF in the configuration.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		generateReport(config.ListEtfs())
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}

func generateReport(etfs []string) {
	type result struct {
		symbol string
		report *report.EtfReport
		err    error
	}
	ch := make(chan result, len(etfs))
	wg := sync.WaitGroup{}
	fmt.Printf("Scraping data...\n")
	for _, s := range etfs {
		wg.Add(1)
		reportGenerator := report.NewReportGenerator(report.WithIntervals([]int{5, 30, 60}))
		go func(symbol string) {
			defer wg.Done()
			etf, err := scraper.ScrapeEtf(symbol)

			var r report.EtfReport
			if err == nil {
				r = *reportGenerator.GenerateReport(*etf)
			}
			res := result{symbol: symbol, report: &r, err: err}
			ch <- res
		}(s)
	}
	reports := make([]report.EtfReport, 0, len(etfs))
	wg.Wait()
	close(ch)

	for r := range ch {
		if r.err != nil {
			fmt.Printf("[%s] Something went wrong while scraping the data. Error details: %s\n", r.symbol, r.err)
			continue
		}
		reports = append(reports, *r.report)
		fmt.Printf("[%s] Success!\n", r.symbol)
		fmt.Printf("%+v\n", *r.report)
	}
	fmt.Printf("Scraping complete.\n")
}
