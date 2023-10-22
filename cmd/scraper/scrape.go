package scraper

import (
	"errors"
	"fmt"
	"strconv"

	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseUrl   = "https://finance.yahoo.com"
	userAgent = "github.com/diogosilva96/etf-cli"
)

var client = &http.Client{}

// EtfExists checks whether the named etfSymbol exists.
func EtfExists(etfSymbol string) bool {
	document, err := getDocument(etfSymbol, client)
	if err != nil {
		return false
	}

	return etfExists(document, etfSymbol)
}

// ScrapesEtf scrapes the ETF for the named etfSymbol and returns the scraped ETF.
func ScrapeEtf(etfSymbol string) (*Etf, error) {
	document, err := getDocument(etfSymbol, client)
	if err != nil {
		return nil, err
	}

	etf, err := scrape(document, etfSymbol)
	if err != nil {
		return nil, err
	}

	return etf, err
}

func createRequest(symbol string) (*http.Request, error) {
	u := buildUrl(symbol)
	r, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	uParsed, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	r.Host = uParsed.Host
	r.Header.Add("User-Agent", userAgent)
	return r, nil
}

func getDocument(symbol string, client *http.Client) (*goquery.Document, error) {
	req, err := createRequest(symbol)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)

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

func scrape(document *goquery.Document, symbol string) (*Etf, error) {
	etf := Etf{}
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

func scrapeHistory(document *goquery.Document, symbol string) ([]EtfHistory, error) {
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

	history := make([]EtfHistory, 0, 60)
	tableData := document.Find("table[data-test=\"historical-prices\"] > tbody")
	tableData.Children().Each(func(row int, rowSelection *goquery.Selection) {
		if row > 60 {
			return
		}
		h := EtfHistory{}
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

func buildUrl(symbol string) string {
	return fmt.Sprintf("%s/quote/%s/history", baseUrl, symbol)
}
