package scraper

// Represents an ETF
type Etf struct {
	symbol, price string
	history       []EtfHistory
}

// Represents a single history element of an ETF
type EtfHistory struct {
	date, price string
}
