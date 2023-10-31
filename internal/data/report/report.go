package report

import (
	"fmt"
	"strings"
)

// EtfReport repesents a report containg etf data.
type EtfReport struct {
	Symbol                              string
	CurrentPrice, Change, PercentChange float32
	IntervalReports                     []EtfIntervalReport
}

// EtfIntervalReport represents a report containing etf for a specific interval.
type EtfIntervalReport struct {
	IntervalInDays                                            int
	MinPrice, MaxPrice, IntervalChange, IntervalPercentChange float32
}

// String outputs a string representation for the report.
func (r *EtfReport) String() string {
	var sb strings.Builder
	if r.Change > 0 {
		sb.WriteString(fmt.Sprintf("[%s] Price: %.2f, Change: +%.2f (+%.2f%%)", r.Symbol, r.CurrentPrice, r.Change, r.PercentChange))
	} else {
		sb.WriteString(fmt.Sprintf("[%s] Price: %.2f, Change: %.2f (%.2f%%)", r.Symbol, r.CurrentPrice, r.Change, r.PercentChange))
	}
	for _, i := range r.IntervalReports {
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("[%v days] Price range: [%.2f, %.2f], Change: %.2f ", i.IntervalInDays, i.MinPrice, i.MaxPrice, i.IntervalChange))
		if i.IntervalPercentChange >= 0 {
			sb.WriteString(fmt.Sprintf("(+%.2f%%)", i.IntervalPercentChange))
			continue
		}
		sb.WriteString(fmt.Sprintf("(%.2f%%)", i.IntervalPercentChange))
	}
	return sb.String()
}
