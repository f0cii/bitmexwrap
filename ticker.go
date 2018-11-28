package bitmexwrap

import (
	"time"
)

type Ticker struct {
	CurrencyPair string
	Last         float64
	High         float64
	Low          float64
	Bid          float64
	Ask          float64
	Volume       float64
	Timestamp    time.Time
}
