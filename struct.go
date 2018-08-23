package coinex

import (
	"time"
)

type Symbol struct {
	Name          string
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	PriceDecimal  int    `json:"price_decimal"`
	AmountDecimal int    `json:"amount_decimal"`
}

type Ticker struct {
	CurrencyPair string
	Last         float64
	High         float64
	Low          float64
	Bid          float64
	Ask          float64
	Volume       float64
}
type DepthInfo struct {
	Price  float64
	Amount float64
}

type Depth struct {
	Sells      []DepthInfo
	Buys       []DepthInfo
	UpdateTime time.Time
}

type Balance struct {
	Currency  string
	Available float64
	Frozen    float64
	Balance   float64
}

type Trade struct {
	Price         float64
	Size          float64
	Side          string
	TickDirection string
	Time          time.Time
}
