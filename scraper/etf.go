package scraper

// Represents an Etf
type Etf struct {
	Symbol  string
	Price   float32
	History []EtfHistory
}

// Represents a single history element of an ETF
type EtfHistory struct {
	Date  string
	Price float32
}
