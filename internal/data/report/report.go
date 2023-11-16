package report

import (
	"fmt"
	"github.com/diogosilva96/etf-cli/internal/data"
	"strings"
)

// EtfReport repesents a report containg etf data.
type EtfReport struct {
	Symbol                              string
	CurrentPrice, Change, PercentChange float32
	IntervalReports                     []EtfIntervalReport
	RawData                             data.Etf
}

// EtfIntervalReport represents a report containing etf for a specific interval.
type EtfIntervalReport struct {
	IntervalInDays                                                                                   int
	MinPrice, MinPriceChange, MinPricePercentChange, MaxPrice, MaxPriceChange, MaxPricePercentChange float32
}

// String outputs a string representation for the report.
func (r EtfReport) String() string {
	var sb strings.Builder
	if r.Change > 0 {
		sb.WriteString(fmt.Sprintf("[%s] Current Price: %.2f (+%.2f)", r.Symbol, r.CurrentPrice, r.Change))
	} else {
		sb.WriteString(fmt.Sprintf("[%s] Current Price: %.2f (%.2f)", r.Symbol, r.CurrentPrice, r.Change))
	}
	for _, i := range r.IntervalReports {
		sb.WriteString("\n")
		if i.MinPriceChange >= 0 {
			sb.WriteString(fmt.Sprintf("[%v days] Min: %.2f (+%.2f)", i.IntervalInDays, i.MinPrice, i.MinPriceChange))
		} else {
			sb.WriteString(fmt.Sprintf("[%v days] Min: %.2f (%.2f)", i.IntervalInDays, i.MinPrice, i.MinPriceChange))
		}
		if i.MaxPriceChange >= 0 {
			sb.WriteString(fmt.Sprintf(", Max: %.2f (+%.2f)", i.MaxPrice, i.MaxPriceChange))
		} else {
			sb.WriteString(fmt.Sprintf(", Max: %.2f (%.2f)", i.MaxPrice, i.MaxPriceChange))
		}
	}
	return sb.String()
}
