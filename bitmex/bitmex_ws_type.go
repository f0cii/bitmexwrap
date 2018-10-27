package bitmex

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	. "github.com/sumorf/coinex"

	log "github.com/sirupsen/logrus"
	"github.com/sumorf/coinex/bitmex/models"
	. "github.com/sumorf/trademodel"
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

type TradeBitmex models.Trade

// Trade Individual & Bucketed Trades
// type TradeBitmex struct {
// 	ForeignNotional float64 `json:"foreignNotional"`
// 	GrossValue      int64   `json:"grossValue"`
// 	HomeNotional    float64 `json:"homeNotional"`
// 	Price           float64 `json:"price"`
// 	Side            string  `json:"side"`
// 	Size            int64   `json:"size"`
// 	Symbol          string  `json:"symbol"`
// 	TickDirection   string  `json:"tickDirection"`
// 	Timestamp       string  `json:"timestamp"`
// 	TrdMatchID      string  `json:"trdMatchID"`
// }

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
	return strconv.FormatInt(o.ID, 10)
}

type OrderBookData []*OrderBookL2

func (od *OrderBookData) GetMap() (ret OrderBookMap) {
	ret = NewOrderBookMap()
	od.GetDataToMap(ret)
	return
}

func (od *OrderBookData) GetDataToMap(ret OrderBookMap) {
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
		case BitmexWSOrderbookL2:
			var orderbooks OrderBookData
			err = json.Unmarshal([]byte(raw), &orderbooks)
			if err != nil {
				return
			}
			r.data = orderbooks
		case BitmexWSTrade:
			var trades []*models.Trade
			err = json.Unmarshal([]byte(raw), &trades)
			if err != nil {
				return
			}
			r.data = trades
		case BitmexWSAnnouncement:
			var announs []Announcement
			err = json.Unmarshal([]byte(raw), &announs)
			if err != nil {
				return
			}
			r.data = announs
		case BitmexWSPosition:
			var pos []*models.Position
			err = json.Unmarshal([]byte(raw), &pos)
			if err != nil {
				return
			}
			r.data = pos
		case BitmexWSOrder:
			var orders []*models.Order
			err = json.Unmarshal([]byte(raw), &orders)
			if err != nil {
				return
			}
			r.data = orders
		case BitmexWSTradeBin1m, BitmexWSTradeBin5m, BitmexWSTradeBin1h, BitmexWSTradeBin1d:
			var klines []*models.TradeBin
			err = json.Unmarshal([]byte(raw), &klines)
			if err != nil {
				return
			}
			r.data = klines
		default:
			log.Debug("unsupport table:", r.Table)
		}
	} else {
	}
	return
}

func (r *Resp) GetTradeData() (trades []*models.Trade) {
	if r.Table != BitmexWSTrade || r.data == nil {
		return
	}
	trades, _ = r.data.([]*models.Trade)
	return
}

func (r *Resp) GetOrderbookL2() (orderbook OrderBookData) {
	if r.Table != BitmexWSOrderbookL2 || r.data == nil {
		return
	}
	orderbook, _ = r.data.(OrderBookData)
	return
}

func (r *Resp) GetPostions() (positions []*models.Position) {
	if r.Table != BitmexWSPosition || r.data == nil {
		return
	}
	positions, _ = r.data.([]*models.Position)
	return
}

func (r *Resp) GetOrders() (orders []*models.Order) {
	if r.Table != BitmexWSOrder || r.data == nil {
		return
	}
	orders, _ = r.data.([]*models.Order)
	return
}

func (r *Resp) GetTradeBin() (klines []*models.TradeBin) {
	if (r.Table != BitmexWSTradeBin1d && r.Table != BitmexWSTradeBin1h && r.Table != BitmexWSTradeBin5m && r.Table != BitmexWSTradeBin1m) || r.data == nil {
		return
	}
	klines, _ = r.data.([]*models.TradeBin)
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
	for _, v := range o {
		if v.Side == "Buy" {
			depth.Buys = append(depth.Buys, DepthInfo{Price: v.Price, Amount: float64(v.Size)})
		} else {
			depth.Sells = append(depth.Sells, DepthInfo{Price: v.Price, Amount: float64(v.Size)})
		}
	}

	sort.Slice(depth.Buys, func(i, j int) bool {
		return depth.Buys[i].Price > depth.Buys[j].Price
	})
	sort.Slice(depth.Sells, func(i, j int) bool {
		return depth.Sells[i].Price < depth.Sells[j].Price
	})
	depth.UpdateTime = time.Now()
	return
}

func transTrade(v *models.Trade) (t Trade) {
	// timestamp: 2018-08-01T05:58:00.737Z
	t = Trade{ID: v.TrdMatchID, Price: v.Price,
		Amount: float64(v.Size), Side: v.Side,
		Remark: v.TickDirection,
		Time:   time.Time(*v.Timestamp)}
	return
}

func transCandle(klines []*models.TradeBin, candles *[]*Candle, binSize string) {
	var secDuration int64
	switch binSize {
	case "1m":
		secDuration = 60
	case "5m":
		secDuration = 60 * 5
	case "1h":
		secDuration = 60 * 60
	case "1d":
		secDuration = 60 * 60 * 24
	}
	for _, v := range klines {
		// v.Timestamp is the close time,not start time
		*candles = append(*candles, &Candle{Start: time.Time(*v.Timestamp).Unix() - secDuration,
			Open:   v.Open,
			High:   v.High,
			Low:    v.Low,
			Close:  v.Close,
			Volume: float64(v.Volume),
			VWP:    v.Vwap,
			Trades: v.Trades})
	}
}

func transOneCandle(v *models.TradeBin) (candle *Candle) {
	candle = &Candle{Start: time.Time(*v.Timestamp).Unix(),
		Open:   v.Open,
		High:   v.High,
		Low:    v.Low,
		Close:  v.Close,
		Volume: float64(v.Volume),
		VWP:    v.Vwap,
		Trades: v.Trades}
	return
}

func transPosition(v *models.Position) (pos *Position) {
	var orderType int
	if v.CurrentQty > 0 {
		orderType = Long
	} else {
		orderType = Short
	}
	if v.CurrentQty == 0 {
		return
	}
	pos = &Position{Info: Contract{Symbol: *v.Symbol, Name: *v.Symbol},
		Type:        orderType,
		Hold:        float64(v.CurrentQty),
		ProfitRatio: float64(v.UnrealisedRoePcnt),
		Price:       v.AvgEntryPrice,
		CostPrice:   v.AvgCostPrice,
		LastPrice:   v.LastPrice,
		MarkPrice:   v.MarkPrice,
	}
	return
}

type PositionMap map[string]*models.Position

func NewPositionMap() (o PositionMap) {
	o = make(PositionMap)
	return
}

func (o PositionMap) Update(pos []*models.Position) {
	var old *models.Position
	var ok bool
	for _, v := range pos {
		if v.CurrentQty == 0 {
			delete(o, *v.Symbol)
		} else {
			old, ok = o[*v.Symbol]
			if ok {
				if v.AvgEntryPrice != 0 {
					old.AvgEntryPrice = v.AvgEntryPrice
				}
				old.CurrentQty = v.CurrentQty
				old.LastPrice = v.LastPrice
				old.MarkPrice = v.MarkPrice
				old.SimpleQty = v.SimpleQty
			} else {
				o[*v.Symbol] = v
			}
		}
	}
	return
}

func (o PositionMap) Pos() (poses []Position) {
	var pos *Position
	for _, v := range o {
		pos = transPosition(v)
		if pos == nil {
			continue
		}
		poses = append(poses, *pos)
	}
	return
}

//type OrderMap map[string]*models.Order

type OrderMap struct {
	m map[string]*models.Order
	sync.RWMutex
}

func NewOrderMap() (o *OrderMap) {
	o = &OrderMap{
		m: make(map[string]*models.Order),
	}
	return
}

func (o *OrderMap) Update(orders []*models.Order, isDelete bool) {
	o.Lock()
	defer o.Unlock()

	var old *models.Order
	var ok bool
	for _, v := range orders {
		if isDelete {
			delete(o.m, *v.OrderID)
		} else {
			old, ok = o.m[*v.OrderID]
			if !ok {
				o.m[*v.OrderID] = v
				continue
			}

			if v.Price > 0 {
				old.Price = v.Price
			}
			if v.OrderQty > 0 {
				old.OrderQty = v.OrderQty
			}
			if v.OrdStatus != "" {
				old.OrdStatus = v.OrdStatus
			}
			if v.AvgPx > 0 {
				old.AvgPx = v.AvgPx
			}
			if v.CumQty > 0 {
				old.CumQty = v.CumQty
			}
			if v.SimpleCumQty > 0 {
				old.SimpleCumQty = v.SimpleCumQty
			}
			old.Timestamp = v.Timestamp // 2018-10-12T02:33:18.886Z
		}
	}
	return
}

func (o *OrderMap) Orders() (orders []Order) {
	o.RLock()
	defer o.RUnlock()

	var pos *Order
	for _, v := range o.m {
		pos = transOrder(v)
		if pos == nil {
			continue
		}
		orders = append(orders, *pos)
	}
	return
}
