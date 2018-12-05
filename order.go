package bitmexwrap

import "time"

type Order struct {
	OrderID  string
	Currency string
	Amount   float64
	Price    float64
	PriceAvg float64
	Status   string
	Side     string
	Type     string
	Time     time.Time
}
