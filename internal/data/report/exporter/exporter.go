package exporter

import "github.com/diogosilva96/etf-cli/internal/data/report"

// ReportExporter represents a component that exports reports.
type ReportExporter interface {
	Export(reports []report.EtfReport) error
}
