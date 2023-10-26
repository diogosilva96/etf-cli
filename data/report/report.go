package report

import (
	"fmt"
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

// String outputs a string for the report.
func (r *EtfReport) String() string {
	s := fmt.Sprintf(`[%s] Price: %.2f, Change: %.2f (%.2f%%)`, r.Symbol, r.CurrentPrice, r.Change, r.PercentChange)
	for _, i := range r.IntervalReports {
		s += fmt.Sprintf("\n")
		s += fmt.Sprintf("[%v days] Price range: [%.2f, %.2f], Change: %.2f ", i.IntervalInDays, i.MinPrice, i.MaxPrice, i.IntervalChange)
		if i.IntervalPercentChange >= 0 {
			s += fmt.Sprintf("(+%.2f%%)", i.IntervalPercentChange)
			continue
		}
		s += fmt.Sprintf("(%.2f%%)", i.IntervalPercentChange)
	}
	return s
}
