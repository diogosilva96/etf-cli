package data

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	baseUrl   = "https://finance.yahoo.com"
	userAgent = "github.com/diogosilva96/etf-cli"
)

// EtfClient Represents an etf client.
type EtfClient struct {
	*http.Client
	BaseUrl   string
	UserAgent string
}

// NewEtfClient creates a new etf client.
func NewEtfClient() *EtfClient {
	return &EtfClient{&http.Client{}, baseUrl, userAgent}
}

// GetEtfDataResponse retrieves the http response for the etf data.
func (c *EtfClient) GetEtfDataResponse(etfSymbol string) (*http.Response, error) {
	if len(strings.TrimSpace(etfSymbol)) == 0 {
		return nil, errors.New("The etf symbol should be specified.")
	}

	req, err := c.createGetEtfRequest(etfSymbol)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

func (c *EtfClient) createGetEtfRequest(etfSymbol string) (*http.Request, error) {
	u := fmt.Sprintf("%s/quote/%s/history", c.BaseUrl, etfSymbol)
	r, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	uParsed, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	r.Host = uParsed.Host
	r.Header.Add("User-Agent", c.UserAgent)
	return r, nil
}
