package bitmex

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/sumorf/coinex/bitmex/models"
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
	order1, err := api.OpenLong(low, 1, false, "")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("OpenLong:", order1)

	order2, err := api.CloseLong(high, 1, false, "")
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
	order1, err := api.OpenShort(high, 2, false, "")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("OpenShort:", order1)
	order2, err := api.CloseShort(low, 2, false, "")
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

func TestTransactionUnmarshal(t *testing.T) {
	s := `[{"transactID":"00000000-0000-0000-0000-000000000000","account":149029,"currency":"XBt","transactType":"UnrealisedPNL","amount":121184,"fee":0,"transactStatus":"Pending","address":"XBTUSD","tx":"","text":"","transactTime":null,"walletBalance":994622,"marginBalance":1115806,"timestamp":null},{"transactID":"00000000-0000-0000-0000-000000000000","account":149029,"currency":"XBt","transactType":"RealisedPNL","amount":-7673,"fee":0,"transactStatus":"Completed","address":"XBTUSD","tx":"","text":"","transactTime":"2018-11-16T12:00:00.000Z","walletBalance":994622,"marginBalance":null,"timestamp":"2018-11-16T12:00:00.000Z"},{"transactID":"c2e2dca9-9e28-89bd-a823-022da81fbd00","account":149029,"currency":"XBt","transactType":"RealisedPNL","amount":1619,"fee":0,"transactStatus":"Completed","address":"XBTUSD","tx":"9fbc8db7-b3a1-85d1-bead-1727c96ba06e","text":"","transactTime":"2018-11-15T12:00:00.000Z","walletBalance":1002295,"marginBalance":null,"timestamp":"2018-11-15T12:00:00.347Z"},{"transactID":"411f59cd-98cc-b173-1588-6f77c17ccaad","account":149029,"currency":"XBt","transactType":"RealisedPNL","amount":792,"fee":0,"transactStatus":"Completed","address":"XBTUSD","tx":"8eda0f7c-8163-dfbe-579a-72ff84206d51","text":"","transactTime":"2018-11-14T12:00:00.000Z","walletBalance":1000676,"marginBalance":null,"timestamp":"2018-11-14T12:00:00.310Z"},{"transactID":"b00c4091-e82c-0104-fc7b-598615b3c32c","account":149029,"currency":"XBt","transactType":"RealisedPNL","amount":-116,"fee":0,"transactStatus":"Completed","address":"XBTUSD","tx":"0e095ac7-2fe3-4eb9-2cdb-725125c44af0","text":"","transactTime":"2018-11-13T12:00:00.000Z","walletBalance":999884,"marginBalance":null,"timestamp":"2018-11-13T12:00:00.237Z"},{"transactID":"6550dbae-f6db-2621-f416-395baa20148d","account":149029,"currency":"XBt","transactType":"Transfer","amount":1000000,"fee":null,"transactStatus":"Completed","address":"0","tx":"b1ff376d-ae54-e0f4-9a50-7bad626395e1","text":"Signup bonus","transactTime":"2018-11-09T01:45:20.474Z","walletBalance":1000000,"marginBalance":null,"timestamp":"2018-11-09T01:45:20.474Z"}]`
	// 2018-11-16T12:00:00.000Z
	trans := []*models.Transaction{}
	err := json.Unmarshal([]byte(s), &trans)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("trans", trans)
}
