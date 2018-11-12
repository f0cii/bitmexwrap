package bitmex

import (
	"errors"
	"fmt"
	"time"

	"github.com/sumorf/coinex/bitmex/client/order"
	"github.com/sumorf/coinex/bitmex/models"
	. "github.com/sumorf/trademodel"
)

const (
	OrderBuy  = "Buy"
	OrderSell = "Sell"

	OrderTypeLimit     = "Limit"
	OrderTypeMarket    = "Market"
	OrderTypeStop      = "Stop"      // stop lose with market price, must set stopPx
	OrderTypeStopLimit = "StopLimit" // stop lose with limit price, must set stopPx
)

var (
	NoOrderFound = errors.New("no such order")
)

func transOrder(o *models.Order) (ret *Order) {
	ret = &Order{OrderID: *o.OrderID,
		Currency: o.Currency,
		Amount:   float64(o.OrderQty),
		Price:    o.Price,
		PriceAvg: o.AvgPx,
		Status:   o.OrdStatus,
		Side:     o.Side,
		Time:     time.Time(o.Timestamp)}
	return
}

// Buy open long with price
func (b *Bitmex) Buy(price float64, amount float64) (ret *Order, err error) {
	ret, err = b.OpenLong(price, amount, false, "")
	return
}

// Buy open long with price
func (b *Bitmex) Sell(price float64, amount float64) (ret *Order, err error) {
	ret, err = b.OpenShort(price, amount, false, "")
	return
}

// OpenLong open long with price
func (b *Bitmex) OpenLong(price float64, amount float64, postOnly bool, timeInForce string) (ret *Order, err error) {
	comment := "open long with bitmex api"
	side := "Buy"
	orderType := "Limit"
	nAmount := int32(amount)
	newOrder, err := b.createOrder(price, nAmount, side, orderType, comment, postOnly, timeInForce)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// CloseLong close long with price
func (b *Bitmex) CloseLong(price float64, amount float64, postOnly bool, timeInForce string) (ret *Order, err error) {
	comment := "close long with bitmex api"
	nAmount := 0 - int32(amount)
	newOrder, err := b.closeOrder(price, nAmount, OrderSell, OrderTypeLimit, comment, postOnly, timeInForce)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// OpenShort open short with price
func (b *Bitmex) OpenShort(price float64, amount float64, postOnly bool, timeInForce string) (ret *Order, err error) {
	comment := "open short with bitmex api"
	nAmount := 0 - int32(amount)
	newOrder, err := b.createOrder(price, nAmount, OrderSell, OrderTypeLimit, comment, postOnly, timeInForce)
	if err != nil {
		return
	}
	ret = transOrder(newOrder)
	return
}

// CloseShort close short with price
func (b *Bitmex) CloseShort(price float64, amount float64, postOnly bool, timeInForce string) (ret *Order, err error) {
	comment := "close short with bitmex api"
	nAmount := int32(amount)
	newOrder, err := b.closeOrder(price, nAmount, OrderBuy, OrderTypeLimit, comment, postOnly, timeInForce)
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
	newOrder, err := b.createOrder(0, nAmount, OrderBuy, OrderTypeMarket, comment, false, "")
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
	newOrder, err := b.closeOrder(0, nAmount, OrderSell, OrderTypeMarket, comment, false, "")
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
	newOrder, err := b.createOrder(0, nAmount, OrderSell, OrderTypeMarket, comment, false, "")
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
	newOrder, err := b.closeOrder(0, nAmount, OrderBuy, OrderTypeMarket, comment, false, "")
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

func (b *Bitmex) closeOrder(price float64, amount int32, side, orderType, comment string, postOnly bool, timeInForce string) (newOrder *models.Order, err error) {
	execInst := "Close"
	if postOnly {
		execInst += ",ParticipateDoNotInitiate"
	}
	params := order.OrderNewParams{
		Side:     &side,
		Symbol:   b.symbol,
		Text:     &comment,
		OrderQty: &amount,
		OrdType:  &orderType,
		ExecInst: &execInst,
	}
	if price != 0 {
		params.Price = &price
	}
	if timeInForce != "" {
		timeInForceString := timeInForce
		params.TimeInForce = &timeInForceString
	}
	orderInfo, err := b.api.Order.OrderNew(&params, nil)
	if err != nil {
		return
	}
	newOrder = orderInfo.Payload
	return
}

func (b *Bitmex) createOrder(price float64, amount int32, side, orderType, comment string, postOnly bool, timeInForce string) (newOrder *models.Order, err error) {
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
	if timeInForce != "" {
		timeInForceString := timeInForce
		params.TimeInForce = &timeInForceString
	}
	if postOnly {
		execInst := "ParticipateDoNotInitiate"
		params.ExecInst = &execInst
	}
	orderInfo, err := b.api.Order.OrderNew(&params, nil)
	if err != nil {
		return
	}
	newOrder = orderInfo.Payload
	return
}

// CancelOrder with oid
func (b *Bitmex) CancelOrder(oid string) (newOrder *Order, err error) {
	comment := "cancle order with bitmex api"
	params := order.OrderCancelParams{
		OrderID: &oid,
		Text:    &comment,
	}
	orderInfo, err := b.api.Order.OrderCancel(&params, nil)
	if err != nil {
		return
	}
	if len(orderInfo.Payload) == 0 {
		err = NoOrderFound
		return
	}
	newOrder = transOrder(orderInfo.Payload[0])
	return
}

// CancelAllOrders cancel all not filled orders
func (b *Bitmex) CancelAllOrders() (orders []*Order, err error) {
	comment := "cancle all order with bitmex api"
	params := order.OrderCancelAllParams{
		Symbol: &b.symbol,
		Text:   &comment,
	}
	orderInfo, err := b.api.Order.OrderCancelAll(&params, nil)
	if err != nil {
		return
	}
	for _, v := range orderInfo.Payload {
		orders = append(orders, transOrder(v))
	}
	return
}

// Orders get all active orders
func (b *Bitmex) Orders() (orders []Order, err error) {
	if b.enableWS {
		orders = b.wsAPI.GetLastOrders()
		return
	}
	return b.GetOrders()
}

// Order get an active orders
func (b *Bitmex) Order(oid string) (newOrder *Order, err error) {
	if b.enableWS {
		order, err := b.wsAPI.GetLastOrder(oid)
		return &order, err
	}
	return b.GetOrder(oid)
}

// GetOrders get all active orders
func (b *Bitmex) GetOrders() (orders []Order, err error) {
	filters := `{"ordStatus":"New"}`
	params := order.OrderGetOrdersParams{
		Symbol: &b.symbol,
		Filter: &filters,
	}
	orderInfo, err := b.api.Order.OrderGetOrders(&params, nil)
	if err != nil {
		return
	}
	for _, v := range orderInfo.Payload {
		orders = append(orders, *transOrder(v))
	}
	return
}

// GetOrder get an active orders
func (b *Bitmex) GetOrder(oid string) (newOrder *Order, err error) {
	filters := fmt.Sprintf(`{"orderID":"%s"}`, oid)
	params := order.OrderGetOrdersParams{
		Symbol: &b.symbol,
		Filter: &filters,
	}
	var orderInfo *order.OrderGetOrdersOK
	orderInfo, err = b.api.Order.OrderGetOrders(&params, nil)
	if err != nil {
		return
	}
	if len(orderInfo.Payload) == 0 {
		err = NoOrderFound
		return
	}
	newOrder = transOrder(orderInfo.Payload[0])
	return
}

func (b *Bitmex) OrderAmend(oid string, price float64) (newOrder *Order, err error) {
	params := order.OrderAmendParams{
		OrderID: &oid,
		Price:   &price,
	}
	var orderAmendOK *order.OrderAmendOK
	orderAmendOK, err = b.api.Order.OrderAmend(&params, nil)
	if err != nil {
		return
	}
	newOrder = transOrder(orderAmendOK.Payload)
	return
}
