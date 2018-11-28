package bitmexwrap

import "time"

type DepthInfo struct {
	Price  float64
	Amount float64
}

type Depth Orderbook

type Orderbook struct {
	Sells      []DepthInfo
	Buys       []DepthInfo
	UpdateTime time.Time
}
