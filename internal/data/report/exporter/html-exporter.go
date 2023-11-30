package exporter

import (
	"github.com/diogosilva96/etf-cli/internal/data/report"
	"github.com/diogosilva96/etf-cli/internal/utils"
	"html/template"
	"os"
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
			symbols: symbols,
			entries: entries,
		},
	}
}

func createChartEntries(r report.EtfReport, entries []chartEntry) []chartEntry {
	for _, h := range r.History {
		found, entry := findChartEntry(entries, h.Date)
		var prices []float32
		if !found {
			entry = &chartEntry{}
			entry.date = h.Date
			entry.prices = append(prices, h.Price)
			entries = append(entries, *entry)
			continue
		}
		entry.prices = append(entry.prices, h.Price)
		entries = append(entries, *entry)
	}
	return entries
}

func findChartEntry(entries []chartEntry, date time.Time) (bool, *chartEntry) {
	for _, e := range entries {
		if e.date == date {
			return true, &e
		}
	}
	return false, nil
}

type templateData struct {
	Date      time.Time
	Reports   []report.EtfReport
	ChartData chartData
}

type chartData struct {
	symbols []string
	entries []chartEntry
}

type chartEntry struct {
	date   time.Time
	prices []float32
}
