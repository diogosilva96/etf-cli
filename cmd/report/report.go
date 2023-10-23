package report

import (
	"errors"
	"fmt"

	"github.com/diogosilva96/etf-cli/cmd/scraper"
)

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

// GenerateReport generates a report based on the provided etf struct.
func GenerateReport(etf scraper.Etf) *EtfReport {

	// set the number of days for the interval reports (e.g., last 5, 30 & 60 days)
	intervals := []int{5, 30, 60}
	report := &EtfReport{
		Symbol:          etf.Symbol,
		CurrentPrice:    etf.Price,
		Change:          (etf.Price - etf.History[1].Price),
		IntervalReports: make([]EtfIntervalReport, 0),
	}

	for _, interval := range intervals {
		intervalReport, err := generateIntervalReport(etf, interval)
		if err != nil {
			// do something?
			continue
		}
		report.IntervalReports = append(report.IntervalReports, *intervalReport)
	}

	return report
}

func generateIntervalReport(etf scraper.Etf, numberOfDays int) (*EtfIntervalReport, error) {

	currentPrice := etf.History[0].Price

	report := &EtfIntervalReport{
		IntervalInDays: numberOfDays,
	}

	max := currentPrice
	min := currentPrice
	var change float32
	historySize := len(etf.History)
	if historySize-1 < numberOfDays {
		return nil, errors.New(fmt.Sprintf("The etf history is only '%v' days long.", historySize))
	}
	for i, h := range etf.History[:numberOfDays] {
		if h.Price > max {
			max = h.Price
		}
		if h.Price < min {
			min = h.Price
		}
		if i == numberOfDays-1 {
			change = (currentPrice - h.Price) - 1
		}
	}
	report.MinPrice = min
	report.MaxPrice = max
	report.IntervalChange = change

	return report, nil
}
