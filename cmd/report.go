package cmd

import (
	"fmt"
	"log"
	"sync"

	"github.com/diogosilva96/etf-cli/cmd/config"
	"github.com/diogosilva96/etf-cli/data"
	"github.com/diogosilva96/etf-cli/data/report"
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
		generateAndPrintReports(config.ListEtfs())
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}

type result struct {
	symbol string
	report *report.EtfReport
	err    error
}

func generateAndPrintReports(etfs []string) {
	ch := make(chan result, len(etfs))
	wg := &sync.WaitGroup{}
	etfClient := data.NewEtfClient()
	fmt.Printf("Fetching etf data...\n")
	for _, s := range etfs {
		wg.Add(1)
		reportGenerator, err := report.NewReportGenerator(report.WithIntervals([]int{5, 30, 60}))
		if err != nil {
			log.Fatal(err) // this should never happen in theory, unless misconfiguration
		}
		go func(etfSymbol string, wg *sync.WaitGroup, ch chan<- result) {
			defer wg.Done()
			etf, err := etfClient.GetEtf(etfSymbol)

			var r report.EtfReport
			if err == nil {
				r = *reportGenerator.GenerateReport(*etf)
			}
			res := result{symbol: etfSymbol, report: &r, err: err}
			ch <- res
		}(s, wg, ch)
	}

	wg.Wait()
	close(ch)

	printReports(ch)
}

func printReports(ch <-chan result) {
	for r := range ch {
		fmt.Printf("----------------------------------------------------------------------------\n")
		if r.err != nil {
			fmt.Printf("Error: [%s] Something went wrong while fetching the etf data. Error details: %s\n", r.symbol, r.err)
			continue
		}

		fmt.Printf("%s\n", r.report.String())
	}
}
