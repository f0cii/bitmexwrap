package bitmex

import (
	"testing"

	. "github.com/SuperGod/coinex"
	log "github.com/sirupsen/logrus"
)

func GetWSClient() (bm *BitmexWS) {
	configs, err := LoadConfigs()
	if err != nil {
		panic(err.Error())
	}
	key, secret := configs.Get("bitmextest")
	bm = NewBitmexWSTest("XBTUSD", key, secret, configs.Proxy)
	return bm
}

func TestTrade(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	clt := GetWSClient()
	err := clt.Connect()
	if err != nil {
		t.Fatal(err.Error())
	}
	tradeChan := make(chan Trade, 1024)
	clt.SetTradeChan(tradeChan)
	for trade := range tradeChan {
		log.Println(trade)
	}
}
