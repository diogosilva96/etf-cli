package scraper

import (
	"errors"
	"fmt"
	"strconv"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/diogosilva96/etf-cli/data"
)

var client = data.NewEtfClient()

// EtfExists checks whether the named etfSymbol exists.
func EtfExists(etfSymbol string) bool {
	document, err := getDocument(etfSymbol)
	if err != nil {
		return false
	}

	return etfExists(document, etfSymbol)
}

// ScrapesEtf scrapes the ETF for the named etfSymbol and returns the scraped ETF.
func ScrapeEtf(etfSymbol string) (*data.Etf, error) {
	document, err := getDocument(etfSymbol)
	if err != nil {
		return nil, err
	}

	etf, err := scrape(document, etfSymbol)
	if err != nil {
		return nil, err
	}

	return etf, err
}

func getDocument(etfSymbol string) (*goquery.Document, error) {
	resp, err := client.GetEtfDataResponse(etfSymbol)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return document, nil
}

func etfExists(document *goquery.Document, symbol string) bool {
	var priceStr = document.Find(fmt.Sprintf("fin-streamer[data-symbol=\"%s\"][data-field=\"regularMarketPrice\"]", symbol)).Text()

	if len(priceStr) <= 0 {
		return false
	}

	return true
}

func scrape(document *goquery.Document, symbol string) (*data.Etf, error) {
	etf := data.Etf{}
	etf.Symbol = symbol
	var priceStr = document.Find(fmt.Sprintf("fin-streamer[data-symbol=\"%s\"][data-field=\"regularMarketPrice\"]", symbol)).Text()

	if len(priceStr) <= 0 {
		return nil, errors.New(fmt.Sprintf("Could not find etf '%s'", symbol))
	}

	price, err := strconv.ParseFloat(priceStr, 32)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not parse price for etf '%s'. Details: %s", symbol, err))
	}
	etf.Price = float32(price)

	history, err := scrapeHistory(document, symbol)
	if err != nil {
		return nil, err
	}

	etf.History = history
	return &etf, nil
}

func scrapeHistory(document *goquery.Document, symbol string) ([]data.EtfHistory, error) {
	tableHeaders := document.Find(fmt.Sprint("table[data-test=\"historical-prices\"] > thead > tr"))
	dateIndex := -1
	priceIndex := -1

	tableHeaders.Children().Each(func(i int, s *goquery.Selection) {
		if dateIndex != -1 && priceIndex != -1 {
			return
		}
		if strings.EqualFold(s.Text(), "date") {
			dateIndex = i
			return
		}

		if strings.EqualFold(s.Text(), "close*") {
			priceIndex = i
		}
	})

	if dateIndex == -1 || priceIndex == -1 {
		return nil, errors.New(fmt.Sprintf("Could not parse the history of the etf '%s'.", symbol))
	}

	history := make([]data.EtfHistory, 0, 90)
	tableData := document.Find("table[data-test=\"historical-prices\"] > tbody")
	tableData.Children().Each(func(row int, rowSelection *goquery.Selection) {
		if row > 60 {
			return
		}
		h := data.EtfHistory{}
		rowSelection.Children().Each(func(col int, colSelection *goquery.Selection) {

			if col == priceIndex {
				price, err := strconv.ParseFloat(colSelection.Text(), 32)
				if err != nil {
					h.Price = 0
					return
				}
				h.Price = float32(price)
			}

			if col == dateIndex {
				h.Date = colSelection.Text()
			}

		})

		if h.Price > 0 && len(h.Date) > 0 {
			history = append(history, h)
		}
	})

	return history, nil
}
