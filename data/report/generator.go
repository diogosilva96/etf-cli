package report

import (
	"errors"
	"fmt"

	"github.com/diogosilva96/etf-cli/data"
)

var (
	defaultIntervals = []int{7, 30}
)

// Represents a report generator.
type ReportGenerator struct {
	Intervals []int
}

// ReportGeneratorOption represents a report generator option.
type ReportGeneratorOption func(*ReportGenerator)

// GenerateReport generates a report based on the provided etf struct.
func (rg *ReportGenerator) GenerateReport(etf data.Etf) *EtfReport {

	// set the number of days for the interval reports (e.g., last 5, 30 & 60 days)
	report := &EtfReport{
		Symbol:          etf.Symbol,
		CurrentPrice:    etf.Price,
		Change:          (etf.Price - etf.History[1].Price),
		IntervalReports: make([]EtfIntervalReport, 0),
	}

	for _, interval := range rg.Intervals {
		intervalReport, err := generateIntervalReport(etf, interval)
		if err != nil {
			// do something, or is it ok to ignore?
			continue
		}
		report.IntervalReports = append(report.IntervalReports, *intervalReport)
	}

	return report
}

// NewReportGenerator initializes a new report generator.
func NewReportGenerator(options ...ReportGeneratorOption) (*ReportGenerator, error) {
	rg := &ReportGenerator{Intervals: defaultIntervals}

	for _, opt := range options {
		opt(rg)
	}

	err := validateReportGenerator(rg)
	if err != nil {
		return nil, err
	}

	return rg, nil
}

// WithIntervals sets the intervals option on a report generator.
func WithIntervals(intervals []int) ReportGeneratorOption {
	return func(rg *ReportGenerator) {
		rg.Intervals = intervals
	}
}

func validateReportGenerator(rg *ReportGenerator) error {
	if rg.Intervals == nil || len(rg.Intervals) == 0 {
		return errors.New("At least one interval should be specified.")
	}
	for _, i := range rg.Intervals {
		if i < 1 {
			return errors.New("The interval should be greater than 0.")
		}
	}
	return nil
}

func generateIntervalReport(etf data.Etf, numberOfDays int) (*EtfIntervalReport, error) {
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
