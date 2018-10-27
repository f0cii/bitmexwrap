package bitmex

import (
	"log"
	"testing"

	. "github.com/sumorf/trademodel"
)

func GetClient() (bm *Bitmex) {
	return GetClientByName("bitmextest", true)
}

func GetHighLowPrice(t *testing.T, api *Bitmex) (high, low float64) {
	price, err := api.Price()
	if err != nil {
		t.Fatal(err.Error())
	}
	low = price - 50
	high = price + 50
	return
}

func TestContracts(t *testing.T) {
	api := GetClient()
	contracts, err := api.Contracts()
	if err != nil {
		// log.Fatal(err.Error())
	}
	log.Printf("%#v\n", contracts)
}

func TestDepth(t *testing.T) {
	api := GetClient()
	depth, err := api.Depth(10)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("depth:%#v\n", depth)
}

func TestTicker(t *testing.T) {
	api := GetClient()
	ticker, err := api.Ticker()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("ticker:%#v\n", ticker)
}

func TestLever(t *testing.T) {
	api := GetClient()
	err := api.SetLever(10)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestPositions(t *testing.T) {
	api := GetClient()
	pos, err := api.Positions()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(pos)
}

func TestBalances(t *testing.T) {
	api := GetClient()
	bals, err := api.ContractBalances()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(bals)
}

func TestUser(t *testing.T) {
	api := GetClient()
	u, err := api.User()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(u)
}

func TestLongLimit(t *testing.T) {

	api := GetClient()
	high, low := GetHighLowPrice(t, api)
	order1, err := api.OpenLong(low, 1)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("OpenLong:", order1)

	order2, err := api.CloseLong(high, 1)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("CloseLong:", order2)
}

func TestLongMarket(t *testing.T) {
	api := GetClient()
	order1, err := api.OpenLongMarket(100)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("OpenLongMarket:", order1)
	// }

	// func TestLongMarketClose(t *testing.T) {
	// 	api := GetClient()
	order2, err := api.CloseLongMarket(100)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("CloseLongMarket:", order2)
}

func TestShortLimit(t *testing.T) {
	api := GetClient()
	high, low := GetHighLowPrice(t, api)
	order1, err := api.OpenShort(high, 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("OpenShort:", order1)
	order2, err := api.CloseShort(low, 2)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("CloseShort:", order2)
}

func TestShortMarket(t *testing.T) {
	api := GetClient()
	order1, err := api.OpenShortMarket(100)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("OpenShortMarket:", order1)
	order2, err := api.CloseShortMarket(100)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("CloseShortMarket:", order2)
}

func TestStopLimit(t *testing.T) {
	api := GetClient()
	high, low := GetHighLowPrice(t, api)
	order1, err := api.StopLoseBuy(high, high+1, 10)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("StopLoseBuy", order1)
	order2, err := api.StopLoseSell(low, low-1, 10)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("StopLoseSell", order2)
}

func TestStopMarket(t *testing.T) {
	api := GetClient()
	order1, err := api.StopLoseBuyMarket(7000, 10)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("StopLoseBuyMarket", order1)
	order2, err := api.StopLoseSellMarket(3000, 10)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("StopLoseBuyMarket", order2)
}

func TestOrders(t *testing.T) {
	api := GetClient()
	orders, err := api.Orders()
	if err != nil {
		t.Fatal(err.Error())
	}
	var order *Order
	t.Log("Orders", orders)
	if len(orders) > 0 {
		order, err = api.CancelOrder(orders[0].OrderID)
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log("cancel Order:", order)
	}
}

func TestCancelAllOrders(t *testing.T) {
	api := GetClient()
	orders, err := api.CancelAllOrders()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("cancel Orders", orders)
}

func TestKlineRecent(t *testing.T) {
	api := GetClient()
	klines, err := api.KlineRecent(10, "1m")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("klines", klines)
}
