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

// Export exports the reports to a html file with the specified file path.
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

	templateData := createHtmlTemplateData(time.Now(), reports)
	err = tmpl.Execute(file, templateData)
	if err != nil {
		return err
	}
	return nil
}

func createHtmlTemplateData(date time.Time, reports []report.EtfReport) htmlTemplateData {
	return htmlTemplateData{Date: date, Reports: reports}
}

type htmlTemplateData struct {
	Date    time.Time
	Reports []report.EtfReport
}
