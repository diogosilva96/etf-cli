package data

import "time"

// Etf Represents an Etf
type Etf struct {
	Symbol  string
	Price   float32
	History []EtfHistory
}

// EtfHistory represents a single history element of an ETF
type EtfHistory struct {
	Date  time.Time
	Price float32
}
