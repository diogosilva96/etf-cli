package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Etf struct {
	symbol, price string
	history       []EtfHistory
}

type EtfHistory struct {
	date, price string
}

const baseUrl = "https://finance.yahoo.com"

func main() {

	// TODO: make this a cli tool and allow removing/adding symbols
	symbols := []string{"VWCE.DE", "asd"}

	etfs := getEtfData(symbols)

	fmt.Printf("%+v\n", etfs)

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
	r.Header.Add("User-Agent", "github.com/diogosilva96/etf-scraper")
	return r, nil
}

func getEtfData(symbols []string) []Etf {

	client := &http.Client{}
	etfs := make([]Etf, 0, len(symbols))

	for _, symbol := range symbols {
		document, err := getDocument(symbol, client)
		if err != nil {
			fmt.Printf("Something went wrong when fetching the data for symbol '%s'. Error details: %s", symbol, err)
		}

		etf, err := scrapeEtf(document, symbol)
		if err != nil {
			fmt.Printf("Something went wrong when scraping the data for symbol '%s'. Error details: %s", symbol, err)
			continue
		}
		etfs = append(etfs, *etf)
	}

	return etfs
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

func scrapeEtf(document *goquery.Document, symbol string) (*Etf, error) {
	etf := Etf{}
	etf.symbol = symbol
	var price = document.Find(fmt.Sprintf("fin-streamer[data-symbol=\"%s\"][data-test=\"qsp-price\"]", symbol)).Text()
	if len(price) <= 0 {
		return nil, errors.New(fmt.Sprintf("Could not find etf with symbol '%s'", symbol))
	}
	etf.price = price

	history, err := scrapeEtfHistory(document)
	if err != nil {
		return nil, err
	}
	etf.history = history

	return &etf, nil
}

func scrapeEtfHistory(document *goquery.Document) ([]EtfHistory, error) {
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

	history := make([]EtfHistory, 0, 60)
	tableData := document.Find("table[data-test=\"historical-prices\"] > tbody")
	tableData.Children().Each(func(row int, rowSelection *goquery.Selection) {
		if row > 60 {
			return
		}
		h := EtfHistory{}
		rowSelection.Children().Each(func(col int, colSelection *goquery.Selection) {

			if col == priceIndex {
				h.price = colSelection.Text()
			}

			if col == dateIndex {
				h.date = colSelection.Text()
			}

		})
		if len(h.price) > 0 && len(h.date) > 0 {
			history = append(history, h)
		}
	})

	return history, nil
}

func buildUrl(symbol string) string {
	return fmt.Sprintf("%s/quote/%s/history", baseUrl, symbol)
}
