package report

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
