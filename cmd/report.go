package cmd

import (
	"log"
	"sync"

	"github.com/diogosilva96/etf-cli/internal/config"
	"github.com/diogosilva96/etf-cli/internal/data"
	"github.com/diogosilva96/etf-cli/internal/data/report"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Provides a report containing real time information of the ETFs in the configuration.",
	Long: `Provides a report containing real time information of the ETFs in the configuration.
	
A report will be generated for each ETF in the configuration.`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		etfs := config.ListEtfs()
		if len(etfs) == 0 {
			cmd.Println("There are no etfs in the configuration.")
			return
		}

		reportGenerator, err := report.NewReportGenerator(report.WithIntervals([]int{7, 30, 60, 90}))
		if err != nil {
			log.Fatal(err) // this should never happen in theory, unless misconfiguration
		}
		ch := make(chan result, len(etfs))
		wg := &sync.WaitGroup{}
		etfClient := data.NewEtfClient()
		cmd.Printf("Retrieving ETF data...\n")
		for _, s := range etfs {
			wg.Add(1)
			go func(etfSymbol string, wg *sync.WaitGroup, ch chan<- result, rg *report.ReportGenerator) {
				defer wg.Done()
				etf, err := etfClient.GetEtf(etfSymbol)

				var r report.EtfReport
				var e error
				if err == nil {
					r, e = rg.GenerateReport(*etf)
					if e != nil {
						cmd.Printf("[%s]\n%s", etfSymbol, e)
					}
				}
				res := result{symbol: etfSymbol, report: &r, data: etf, err: err}
				ch <- res
			}(s, wg, ch, reportGenerator)
		}

		wg.Wait()
		close(ch)

		printReports(cmd, ch)
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}

type result struct {
	symbol string
	report *report.EtfReport
	data   *data.Etf
	err    error
}

func printReports(cmd *cobra.Command, ch <-chan result) {
	for r := range ch {
		cmd.Printf("----------------------------------------------------------------------------\n")
		if r.err != nil {
			cmd.Printf("[%s]\n", r.symbol)
			cmd.Printf("Something went wrong: %s\n", r.err)
			continue
		}

		cmd.Printf("%s\n", r.report.String())
	}
}
