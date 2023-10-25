package report

import "fmt"

// EtfReport repesents a report containg etf data.
type EtfReport struct {
	Symbol               string
	CurrentPrice, Change float32
	IntervalReports      []EtfIntervalReport
}

// EtfIntervalReport represents a report containing etf for a specific interval.
type EtfIntervalReport struct {
	IntervalInDays                     int
	MinPrice, MaxPrice, IntervalChange float32
}

//String outputs a string for the report.
func (r *EtfReport) String() string {
	s := fmt.Sprintf(`[%s] Price: %v, Change: %v`, r.Symbol, r.CurrentPrice, r.Change)
	for _, i := range r.IntervalReports {
		s += fmt.Sprintf("\n")
		s += fmt.Sprintf(`[%v days] Price range: [%v, %v], Change: %v,`, i.IntervalInDays, i.MinPrice, i.MaxPrice, i.IntervalChange)
	}
	return s
}
