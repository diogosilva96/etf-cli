package exporter

import (
	"github.com/diogosilva96/etf-cli/internal/data/report"
	"github.com/diogosilva96/etf-cli/internal/utils"
	"html/template"
	"os"
	"sort"
	"time"
)

// HtmlReportExporter represents an html exporter component.
type htmlReportExporter struct {
}

// NewHtmlReportExporter creates a new html report exporter with the specified options.
func NewHtmlReportExporter() ReportExporter {
	return htmlReportExporter{}
}

// Export exports the Reports to a html file with the specified file path.
func (e htmlReportExporter) Export(reports []report.EtfReport, filePath string) error {
	err := utils.CreateFoldersIfNotExist(filePath)
	if err != nil {
		return err
	}

	tmpl := template.Must(template.ParseFiles("./templates/report.html"))

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	defer file.Close()
	if err != nil {
		return err
	}

	templateData := createTemplateData(time.Now(), reports)
	err = tmpl.Execute(file, templateData)
	if err != nil {
		return err
	}
	return nil
}

func createTemplateData(date time.Time, reports []report.EtfReport) templateData {
	var symbols []string
	var entries []chartEntry
	for _, r := range reports {
		symbols = append(symbols, r.Symbol)
		entries = createChartEntries(r, entries)
	}
	return templateData{
		Date:    date,
		Reports: reports,
		ChartData: chartData{
			Symbols: symbols,
			Entries: entries,
		},
	}
}

func createChartEntries(r report.EtfReport, entries []chartEntry) []chartEntry {
	for _, h := range r.History {
		found, idx := findChartEntryIndexByDate(entries, h.Date)
		var prices []float32
		if found {
			entry := entries[idx]
			entry.Prices = append(entry.Prices, h.Price)
			entries = removeEntry(entries, idx) // remove "old" found entry
			entries = append(entries, entry)    // append "new" entry
			continue
		}
		entries = append(entries, chartEntry{
			Date:   h.Date,
			Prices: append(prices, h.Price),
		})
	}

	sortByDateAscending(entries)

	return entries
}

func removeEntry(entries []chartEntry, idx int) []chartEntry {
	return append(entries[:idx], entries[idx+1:]...)
}

func sortByDateAscending(entries []chartEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Date.Before(entries[j].Date)
	})
}

func findChartEntryIndexByDate(entries []chartEntry, date time.Time) (bool, int) {
	for idx, e := range entries {
		if e.Date.Year() == date.Year() && e.Date.YearDay() == date.YearDay() {
			return true, idx
		}
	}
	return false, -1
}

type templateData struct {
	Date      time.Time
	Reports   []report.EtfReport
	ChartData chartData
}

type chartData struct {
	Symbols []string
	Entries []chartEntry
}

type chartEntry struct {
	Date   time.Time
	Prices []float32
}
