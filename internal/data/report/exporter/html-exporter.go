package exporter

import (
	"fmt"
	"github.com/diogosilva96/etf-cli/internal/data/report"
	"html/template"
	"os"
	"time"
)

const (
	// timestampFormat represents a time format for 'yyyyMMddhhmmss'.
	timestampFormat = "20060102150405"
)

// HtmlReportExporter represents an html exporter component.
type htmlReportExporter struct {
	fileNameFunc func() string
}

// NewHtmlReportExporter creates a new html report exporter with the specified options.
func NewHtmlReportExporter(options ...func(*htmlReportExporter)) ReportExporter {
	exporter := htmlReportExporter{}
	for _, opt := range options {
		opt(&exporter)
	}

	if exporter.fileNameFunc == nil {
		exporter.fileNameFunc = func() string {
			return fmt.Sprintf("%s-report.%s", time.Now().Format(timestampFormat), "html")
		}
	}

	return exporter
}

// WithFileNameFunc provides an option to configure the file name func for the exported html report.
func WithFileNameFunc(fileNameFunc func() string) func(*htmlReportExporter) {
	return func(e *htmlReportExporter) {
		e.fileNameFunc = fileNameFunc
	}
}

// Export exports the reports to a html file with the specified file name.
func (e htmlReportExporter) Export(reports []report.EtfReport) error {
	// TODO: make output path configurable
	outPath := "./out"
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		err = os.Mkdir(outPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	fileName := e.fileNameFunc()
	filePath := fmt.Sprintf("%s/%s", outPath, fileName)
	tmpl := template.Must(template.ParseFiles("./templates/report.html"))

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		return err
	}

	rt := reportTemplate{Date: time.Now(), Reports: reports}
	err = tmpl.Execute(f, rt)
	return err
}

type reportTemplate struct {
	Date    time.Time
	Reports []report.EtfReport
}

func createReportFileName(t time.Time) string {
	// TODO: move this to somewhere else
	return fmt.Sprintf("%s-report.%s", t.Format(timestampFormat), "html")
}
