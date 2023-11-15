package report

import (
	"errors"
	"fmt"
	"math"

	"github.com/diogosilva96/etf-cli/internal/data"
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
func (rg ReportGenerator) GenerateReport(etf data.Etf) (EtfReport, error) {
	// set the number of days for the interval reports (e.g., last 5, 30 & 60 days)
	previousDayPrice := etf.History[1].Price
	report := EtfReport{
		Symbol:        etf.Symbol,
		CurrentPrice:  etf.Price,
		Change:        calculateChange(etf.Price, previousDayPrice),
		PercentChange: calculatePercentChange(etf.Price, previousDayPrice),
		RawData:       etf,
	}

	var errs []error
	for _, interval := range rg.Intervals {
		intervalReport, err := generateIntervalReport(etf, interval)
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Something went wrong while generating report for %v days interval: %s", interval, err)))
			continue
		}
		report.IntervalReports = append(report.IntervalReports, *intervalReport)
	}
	return report, errorsToError(errs)
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

	for _, h := range etf.History[:numberOfDays] {
		report.MaxPrice = calculateMax(h.Price, report.MaxPrice)
		report.MinPrice = calculateMin(h.Price, report.MinPrice)
	}

	report.MinPriceChange = calculateChange(etf.Price, report.MinPrice)
	report.MinPricePercentChange = calculatePercentChange(etf.Price, report.MinPrice)
	report.MaxPriceChange = calculateChange(etf.Price, report.MaxPrice)
	report.MaxPricePercentChange = calculatePercentChange(etf.Price, report.MaxPrice)

	return report, nil
}

func calculatePercentChange(currentValue float32, previousValue float32) float32 {
	return ((currentValue - previousValue) / float32(math.Abs(float64(previousValue)))) * 100
}

func calculateChange(currentValue float32, previousValue float32) float32 {
	return currentValue - previousValue
}

func calculateMax(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func calculateMin(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func errorsToError(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	var s string
	for _, e := range errs {
		s += fmt.Sprintf("- %s\n", e)
	}
	return errors.New(s)
}
