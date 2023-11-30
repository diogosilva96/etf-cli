package exporter

import "github.com/diogosilva96/etf-cli/internal/data/report"

// ReportExporter represents a component that exports reports.
type ReportExporter interface {
	// Export exports the specified reports to the specified filePath.
	Export(reports []report.EtfReport, filePath string) error
}
