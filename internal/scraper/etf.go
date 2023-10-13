package scraper

// Represents an Etf
type Etf struct {
	Symbol, Price string
	History       []EtfHistory
}

// Represents a single history element of an ETF
type EtfHistory struct {
	Date, Price string
}
