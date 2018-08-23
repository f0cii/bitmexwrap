package bitmex

import (
	"encoding/json"
	"fmt"
	"time"

	. "github.com/SuperGod/coinex"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Welcome struct {
	Info      string
	Version   string
	Timestamp time.Time
	Docs      string
	Limit     struct {
		Remaining int
	}
}

// OrderBookL2 contains order book l2
type OrderBookL2 struct {
	ID     int64   `json:"id"`
	Price  float64 `json:"price"`
	Side   string  `json:"side"`
	Size   int64   `json:"size"`
	Symbol string  `json:"symbol"`
}

// Trade Individual & Bucketed Trades
type TradeBitmex struct {
	ForeignNotional float64 `json:"foreignNotional"`
	GrossValue      int64   `json:"grossValue"`
	HomeNotional    float64 `json:"homeNotional"`
	Price           float64 `json:"price"`
	Side            string  `json:"side"`
	Size            int64   `json:"size"`
	Symbol          string  `json:"symbol"`
	TickDirection   string  `json:"tickDirection"`
	Timestamp       string  `json:"timestamp"`
	TrdMatchID      string  `json:"trdMatchID"`
}

// Announcement General Announcements
type Announcement struct {
	Content string `json:"content"`
	Date    string `json:"date"`
	ID      int32  `json:"id"`
	Link    string `json:"link"`
	Title   string `json:"title"`
}

type WSCmd struct {
	Command string        `json:"op"`
	Args    []interface{} `json:"args"`
}

type SubscribeResp struct {
	Success   bool   `json:"success"`
	Subscribe string `json:"subscribe"`
}

type ErrorResponse struct {
	Status int         `json:"status"`
	Error  string      `json:"error"`
	Meta   interface{} `json:"meta"`
}

type MainResponse struct {
	Table string   `json:"table"`
	Keys  []string `json:"keys"`
	Types struct {
		ID     string `json:"id"`
		Price  string `json:"price"`
		Side   string `json:"side"`
		Size   string `json:"size"`
		Symbol string `json:"symbol"`
	} `json:"types"`
	ForeignKeys struct {
		Side   string `json:"side"`
		Symbol string `json:"symbol"`
	} `json:"foreignKeys"`
	Attributes struct {
		ID     string `json:"id"`
		Symbol string `json:"symbol"`
	} `json:"Attributes"`
}

func (o *OrderBookL2) Key() string {
	// return fmt.Sprintf("%s-%s-%d", o.Symbol, o.Side, o.ID)
	return string(o.ID)
}

type OrderBookData []*OrderBookL2

func (od *OrderBookData) GetMap() (ret map[string]*OrderBookL2) {
	od.GetDataToMap(ret)
	return
}

func (od *OrderBookData) GetDataToMap(ret map[string]*OrderBookL2) {
	for _, v := range *od {
		ret[v.Key()] = v
	}
	return
}

type Resp struct {
	Request WSCmd `json:"request"`
	SubscribeResp
	ErrorResponse
	MainResponse

	Action string

	data       interface{}
	hasStatus  bool
	hasSuccess bool
	hasTable   bool
}

func (r *Resp) Decode(buf []byte) (err error) {
	err = json.Unmarshal(buf, r)
	if err != nil {
		return
	}
	ret := gjson.ParseBytes(buf)
	if !ret.IsObject() {
		err = fmt.Errorf("unmarshal json error")
		return
	}
	if ret.Get("status").Exists() {
		r.hasStatus = true
	} else if ret.Get("success").Exists() {
		r.hasSuccess = true
	} else if ret.Get("table").Exists() {
		r.hasTable = true
		raw := ret.Get("data").Raw
		switch r.Table {
		case bitmexWSOrderbookL2:
			var orderbooks OrderBookData
			err = json.Unmarshal([]byte(raw), &orderbooks)
			if err != nil {
				return
			}
			r.data = orderbooks
		case bitmexWSTrade:
			var trades []TradeBitmex
			err = json.Unmarshal([]byte(raw), &trades)
			if err != nil {
				return
			}
			r.data = trades
		case bitmexWSAnnouncement:
			var announs []Announcement
			err = json.Unmarshal([]byte(raw), &announs)
			if err != nil {
				return
			}
			r.data = announs
		default:
			log.Debug("unsupport table:", r.Table)
		}
	} else {
	}
	return
}

func (r *Resp) GetTradeData() (trades []TradeBitmex) {
	if r.Table != bitmexWSTrade || r.data == nil {
		return
	}
	trades, _ = r.data.([]TradeBitmex)
	return
}

func (r *Resp) GetOrderbookL2() (orderbook OrderBookData) {
	if r.Table != bitmexWSOrderbookL2 || r.data == nil {
		return
	}
	orderbook, _ = r.data.(OrderBookData)
	return
}

func (r *Resp) HasStatus() bool {
	return r.hasStatus
}

func (r *Resp) HasSuccess() bool {
	return r.hasSuccess
}

func (r *Resp) HasTable() bool {
	return r.hasTable
}

type OrderBookMap map[string]*OrderBookL2

func NewOrderBookMap() (o OrderBookMap) {
	o = make(OrderBookMap)
	return
}

func (o OrderBookMap) GetDepth() (depth Depth) {
	return
}

func transTrade(v TradeBitmex) (t Trade) {
	// timestamp: 2018-08-01T05:58:00.737Z
	tm, err := time.Parse(time.RFC3339, v.Timestamp)
	if err != nil {
		log.Info("parse time failed:", v.Timestamp)
	}
	t = Trade{Price: v.Price, Size: float64(v.Size), Side: v.Side, TickDirection: v.TickDirection, Time: tm}
	return
}
