package bitmex

import (
	"time"

	"github.com/SuperGod/coinex/bitmex/client/order"
	"github.com/SuperGod/coinex/bitmex/models"
	. "github.com/SuperGod/trademodel"
)

const (
	OrderBuy  = "Buy"
	OrderSell = "Sell"

	OrderTypeLimit     = "Limit"
	OrderTypeMarket    = "Market"
	OrderTypeStop      = "Stop"      // stop lose with market price, must set stopPx
	OrderTypeStopLimit = "StopLimit" // stop lose with limit price, must set stopPx
)

func transOrder(o *models.Order) (ret *Order) {
	ret = &Order{OrderID: *o.OrderID,
		Currency: o.Currency,
		Amount:   float64(o.OrderQty),
		Price:    o.AvgPx,
		Status:   o.OrdStatus,
		Side:     o.Side,
		Time:     time.Time(o.Timestamp)}
	return
}

// OpenLong open long with price
func (b *Bitmex) OpenLong(price float64, amount float64) (ret *Order, err error) {
	comment := "open long with bitmex api"
	side := "Buy"
	orderType := "Limit"
	nAmount := int32(amount)
	newOrder, err := b.createOrder(price, nAmount, side, orderType, comment)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// CloseLong close long with price
func (b *Bitmex) CloseLong(price float64, amount float64) (ret *Order, err error) {
	comment := "close long with bitmex api"
	nAmount := 0 - int32(amount)
	newOrder, err := b.createOrder(price, nAmount, OrderSell, OrderTypeLimit, comment)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// OpenShort open short with price
func (b *Bitmex) OpenShort(price float64, amount float64) (ret *Order, err error) {
	comment := "open short with bitmex api"
	nAmount := 0 - int32(amount)
	newOrder, err := b.createOrder(price, nAmount, OrderSell, OrderTypeLimit, comment)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// CloseShort close short with price
func (b *Bitmex) CloseShort(price float64, amount float64) (ret *Order, err error) {
	comment := "close short with bitmex api"
	nAmount := int32(amount)
	newOrder, err := b.createOrder(price, nAmount, OrderBuy, OrderTypeLimit, comment)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// OpenLongMarket open long with market price
func (b *Bitmex) OpenLongMarket(amount float64) (ret *Order, err error) {
	comment := "open market long with bitmex api"
	nAmount := int32(amount)
	newOrder, err := b.createOrder(0, nAmount, OrderBuy, OrderTypeMarket, comment)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// CloseLongarket close long with market price
func (b *Bitmex) CloseLongMarket(amount float64) (ret *Order, err error) {
	comment := "close market long with bitmex api"
	nAmount := 0 - int32(amount)
	newOrder, err := b.createOrder(0, nAmount, OrderSell, OrderTypeMarket, comment)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// OpenShortMarket open short with market price
func (b *Bitmex) OpenShortMarket(amount float64) (ret *Order, err error) {
	comment := "open market short with bitmex api"
	nAmount := 0 - int32(amount)
	newOrder, err := b.createOrder(0, nAmount, OrderSell, OrderTypeMarket, comment)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// CloseShortMarket close short with market price
func (b *Bitmex) CloseShortMarket(amount float64) (ret *Order, err error) {
	comment := "close market short with bitmex api"
	nAmount := int32(amount)
	newOrder, err := b.createOrder(0, nAmount, OrderBuy, OrderTypeMarket, comment)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// StopLoseBuy when marketPrice>=stopPrice, create buy order with price and amount
func (b *Bitmex) StopLoseBuy(stopPrice, price, amount float64) (ret *Order, err error) {
	comment := "stop limit buy with bitmex api"
	nAmount := int32(amount)
	newOrder, err := b.createStopOrder(stopPrice, price, nAmount, OrderBuy, OrderTypeStopLimit, comment)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// StopLoseSell when marketPrice<=stopPrice, create sell order with price and amount
func (b *Bitmex) StopLoseSell(stopPrice, price, amount float64) (ret *Order, err error) {
	comment := "stop limit buy with bitmex api"
	nAmount := int32(amount)
	newOrder, err := b.createStopOrder(stopPrice, price, nAmount, OrderSell, OrderTypeStopLimit, comment)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// StopLoseSell when marketPrice>=stopPrice, create buy order with marketPrice and amount
func (b *Bitmex) StopLoseBuyMarket(price, amount float64) (ret *Order, err error) {
	comment := "stop market buy with bitmex api"
	nAmount := int32(amount)
	newOrder, err := b.createStopOrder(price, 0, nAmount, OrderBuy, OrderTypeStop, comment)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// StopLoseSell when marketPrice<=stopPrice, create buy order with marketPrice and amount
func (b *Bitmex) StopLoseSellMarket(price, amount float64) (ret *Order, err error) {
	comment := "stop market buy with bitmex api"
	nAmount := int32(amount)
	newOrder, err := b.createStopOrder(price, 0, nAmount, OrderSell, OrderTypeStop, comment)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// createStopOrder stop order
// if orderType is Buy, when marketPrice>= price, buy
// if orderType is Sell, when marketPrice<=price, sell
func (b *Bitmex) createStopOrder(stopPrice, price float64, amount int32, side, orderType, comment string) (newOrder *models.Order, err error) {
	params := order.OrderNewParams{
		StopPx:   &stopPrice,
		Side:     &side,
		Symbol:   b.symbol,
		Text:     &comment,
		OrderQty: &amount,
		OrdType:  &orderType,
	}
	if price != 0 {
		params.Price = &price
	}
	orderInfo, err := b.api.Order.OrderNew(&params, nil)
	if err != nil {
		return
	}
	newOrder = orderInfo.Payload
	return
}

func (b *Bitmex) createOrder(price float64, amount int32, side, orderType, comment string) (newOrder *models.Order, err error) {
	params := order.OrderNewParams{
		Side:     &side,
		Symbol:   b.symbol,
		Text:     &comment,
		OrderQty: &amount,
		OrdType:  &orderType,
	}
	if price != 0 {
		params.Price = &price
	}
	orderInfo, err := b.api.Order.OrderNew(&params, nil)
	if err != nil {
		return
	}
	newOrder = orderInfo.Payload
	return
}

// CancelOrder with oid
func (b *Bitmex) CancelOrder(oid string) (orders []*models.Order, err error) {
	comment := "cancle order with bitmex api"
	params := order.OrderCancelParams{
		OrderID: &oid,
		Text:    &comment,
	}
	orderInfo, err := b.api.Order.OrderCancel(&params, nil)
	if err != nil {
		return
	}
	orders = orderInfo.Payload
	return
}

// CancelAllOrders cancel all not filled orders
func (b *Bitmex) CancelAllOrders() (orders []*models.Order, err error) {
	comment := "cancle all order with bitmex api"
	params := order.OrderCancelAllParams{
		Symbol: &b.symbol,
		Text:   &comment,
	}
	orderInfo, err := b.api.Order.OrderCancelAll(&params, nil)
	if err != nil {
		return
	}
	orders = orderInfo.Payload
	return
}

// Orders get all active orders
func (b *Bitmex) Orders() (orders []*models.Order, err error) {
	filters := `{"ordStatus":"New"}`
	params := order.OrderGetOrdersParams{
		Symbol: &b.symbol,
		Filter: &filters,
	}
	orderInfo, err := b.api.Order.OrderGetOrders(&params, nil)
	if err != nil {
		return
	}
	orders = orderInfo.Payload
	return
}
