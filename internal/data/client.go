package data

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	defaultBaseUrl    = "https://finance.yahoo.com"
	defaultUserAgent  = "github.com/diogosilva96/etf-cli"
	host              = "finance.yahoo.com"
	numberHistoryDays = 180
)

// EtfClient Represents an etf client.
type EtfClient struct {
	client    *http.Client
	BaseUrl   string
	UserAgent string
}

// NewEtfClient creates a new etf client.
func NewEtfClient() *EtfClient {
	c := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &headerTransport{
			headers: map[string]string{
				"User-Agent": defaultUserAgent,
				"Host":       host,
			},
		},
	}

	return &EtfClient{c, defaultBaseUrl, defaultUserAgent}
}

// EtfExists checks whether the named etfSymbol exists.
func (c EtfClient) EtfExists(etfSymbol string) bool {
	resp, err := c.getEtfDataResponse(etfSymbol)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return false
	}

	return etfExists(document, etfSymbol)
}

// GetEtf retrieves the ETF information for the named etfSymbol and returns the scraped ETF.
func (c EtfClient) GetEtf(etfSymbol string) (*Etf, error) {
	resp, err := c.getEtfDataResponse(etfSymbol)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	etf, err := scrapeEtf(document, etfSymbol)
	if err != nil {
		return nil, err
	}

	return etf, nil
}

func (c EtfClient) getEtfDataResponse(etfSymbol string) (*http.Response, error) {
	if len(strings.TrimSpace(etfSymbol)) == 0 {
		return nil, errors.New("The etf symbol should be specified.")
	}

	req, err := c.createGetEtfRequest(etfSymbol)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c EtfClient) createGetEtfRequest(etfSymbol string) (*http.Request, error) {
	u := fmt.Sprintf("%s/quote/%s/history", c.BaseUrl, etfSymbol)
	r, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func scrapeEtf(document *goquery.Document, symbol string) (*Etf, error) {
	etf := &Etf{}
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
	return etf, nil
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

	history := make([]EtfHistory, 0, numberHistoryDays)
	tableData := document.Find("table[data-test=\"historical-prices\"] > tbody")
	tableData.Children().Each(func(row int, rowSelection *goquery.Selection) {
		if row > numberHistoryDays {
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

func etfExists(document *goquery.Document, symbol string) bool {
	var priceStr = document.Find(fmt.Sprintf("fin-streamer[data-symbol=\"%s\"][data-field=\"regularMarketPrice\"]", symbol)).Text()

	if len(priceStr) <= 0 {
		return false
	}

	return true
}
