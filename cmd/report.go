package cmd

import (
	"fmt"
	"github.com/diogosilva96/etf-cli/internal/config"
	"github.com/diogosilva96/etf-cli/internal/data"
	"github.com/diogosilva96/etf-cli/internal/data/report"
	"github.com/spf13/cobra"
	"html/template"
	"log"
	"os"
	"sync"
	"time"
)

const (
	// htmlFileType represents an html file type.
	htmlFileType fileType = "html"
)

const (
	// timestampFormat represents a time format for 'yyyyMMddhhmmss'.
	timestampFormat = "20060102150405"
)

// fileType represents a file type
type fileType string

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
				res := result{symbol: etfSymbol, report: &r, err: err}
				ch <- res
			}(s, wg, ch, reportGenerator)
		}

		wg.Wait()
		close(ch)

		reports := getReports(*cmd, ch)
		printReports(*cmd, reports)

		// TODO: move this to somewhere else
		fileName := createReportFileName(time.Now(), htmlFileType)
		err = export(reports, fileName)
		if err != nil {
			cmd.Printf("Something went wrong while exporting report: %s\n", err)
		}
	},
}

func getReports(cmd cobra.Command, ch <-chan result) []report.EtfReport {
	var reports []report.EtfReport
	for r := range ch {

		if r.err != nil {
			cmd.Printf("----------------------------------------------------------------------------\n")
			cmd.Printf("[%s]\n", r.symbol)
			cmd.Printf("Something went wrong: %s\n", r.err)
			continue
		}
		reports = append(reports, *r.report)
	}
	return reports
}

func printReports(cmd cobra.Command, reports []report.EtfReport) {
	for _, r := range reports {
		cmd.Printf("----------------------------------------------------------------------------\n")
		cmd.Printf("%s\n", r.String())
	}
}

func init() {
	rootCmd.AddCommand(reportCmd)
}

type result struct {
	symbol string
	report *report.EtfReport
	err    error
}

func createReportFileName(t time.Time, fType fileType) string {
	// TODO: move this to somewhere else
	return fmt.Sprintf("%s-report.%s", t.Format(timestampFormat), fType)
}

func export(reports []report.EtfReport, fileName string) error {
	// TODO: move this to somewhere else
	outPath := "./out"
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		err = os.Mkdir(outPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	filePath := fmt.Sprintf("%s/%s", outPath, fileName)
	tmpl := template.Must(template.ParseFiles("./templates/report.html"))

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		return err
	}

	type reportTemplate struct {
		Date    time.Time
		Reports []report.EtfReport
	}
	rt := reportTemplate{Date: time.Now(), Reports: reports}
	err = tmpl.Execute(f, rt)
	return err
}
