package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
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

	// TODO: completely remove colly & simply use standard lib http
	collector := colly.NewCollector()
	configure(collector)

	etfs := getEtfData(collector, symbols)

	fmt.Printf("%+v\n", etfs)

}

func configure(collector *colly.Collector) {
	collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request to URL '", r.Request.URL, "' failed with response '", r, "'\nError:", err)
	})
}

func getEtfData(collector *colly.Collector, symbols []string) []Etf {
	etfs := make([]Etf, 0, len(symbols))

	collector.OnResponse(func(r *colly.Response) {
		reader := bytes.NewReader(r.Body)
		document, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			log.Fatal(err)
		}

		symbol := strings.Split(r.Request.URL.Path, "/")[2]
		etf, err := scrapeEtf(document, symbol)
		if err != nil {
			fmt.Println(err)
		} else {
			etfs = append(etfs, *etf)
		}
	})

	visitUrls(symbols, collector)

	return etfs
}

func visitUrls(symbols []string, collector *colly.Collector) {
	for _, symbol := range symbols {
		url := buildUrl(symbol)
		collector.Visit(url)
	}
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
