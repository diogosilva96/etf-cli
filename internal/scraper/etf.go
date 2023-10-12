package scraper

type Etf struct {
	symbol, price string
	history       []EtfHistory
}

type EtfHistory struct {
	date, price string
}
