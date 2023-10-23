package report

import (
	"errors"
	"fmt"

	"github.com/diogosilva96/etf-cli/cmd/scraper"
)

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
			// do something, or is it ok to ignore?
			continue
		}
		report.IntervalReports = append(report.IntervalReports, *intervalReport)
	}

	return report
}

func generateIntervalReport(etf scraper.Etf, numberOfDays int) (*EtfIntervalReport, error) {
	report := &EtfIntervalReport{
		IntervalInDays: numberOfDays,
		MaxPrice:       etf.Price,
		MinPrice:       etf.Price,
	}

	historySize := len(etf.History)
	if historySize-1 < numberOfDays {
		return nil, errors.New(fmt.Sprintf("The etf history is only '%v' days long.", historySize))
	}
	for i, h := range etf.History[:numberOfDays] {
		if h.Price > report.MaxPrice {
			report.MaxPrice = h.Price
		}
		if h.Price < report.MinPrice {
			report.MinPrice = h.Price
		}
		if i == numberOfDays-1 {
			report.IntervalChange = (etf.Price - h.Price) - 1
		}
	}

	return report, nil
}
