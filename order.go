package bitmexwrap

import "time"

type Order struct {
	OrderID    string
	Currency   string
	Amount     float64 // 委托数量
	DealAmount float64 // 成交数量
	Price      float64 // 委托价格
	PriceAvg   float64 // 平均成交价
	Status     string  // 状态
	Side       string
	Type       string
	Time       time.Time
}
