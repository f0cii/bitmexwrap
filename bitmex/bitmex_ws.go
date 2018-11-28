package bitmex

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-openapi/strfmt"

	. "github.com/sumorf/coinex"

	. "github.com/sumorf/trademodel"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/sumorf/coinex/bitmex/models"
)

const (
	MaxTableLen     = 200
	bitmexWSURL     = "wss://www.bitmex.com/realtime"
	testBitmexWSURL = "wss://testnet.bitmex.com/realtime"

	// Bitmex websocket op
	BitmexWSOrderbookL2  = "orderBookL2" // Full level 2 orderBook
	BitmexWSOrderbookL10 = "orderBook10" // Top 10 levels using traditional full book push

	BitmexWSTrade      = "trade"      // Live trades
	BitmexWSTradeBin1m = "tradeBin1m" // 1-minute trade bins
	BitmexWSTradeBin5m = "tradeBin5m" // 5-minute trade bins
	BitmexWSTradeBin1h = "tradeBin1h" // 1-hour trade bins
	BitmexWSTradeBin1d = "tradeBin1d" // 1-day trade bins

	BitmexWSAnnouncement = "announcement" // Site announcements
	BitmexWSLiquidation  = "liquidation"  // Liquidation orders as they're entered into the book

	BitmexWSQuote      = "quote"      // Top level of the book
	BitmexWSQuoteBin1m = "quoteBin1m" // 1-minute quote bins
	BitmexWSQuoteBin5m = "quoteBin5m" // 5-minute quote bins
	BitmexWSQuoteBin1h = "quoteBin1h" // 1-hour quote bins
	BitmexWSQuoteBin1d = "quoteBin1d" // 1-day quote bins

	// Bitmex websocket private op
	BitmexWSExecution = "execution" // Individual executions; can be multiple per order
	BitmexWSOrder     = "order"     // Live updates on your orders
	BitmexWSMargin    = "margin"    // Updates on your current account balance and margin requirements
	BitmexWSPosition  = "position"  // Updates on your positions

	bitmexActionInitialData = "partial"
	bitmexActionInsertData  = "insert"
	bitmexActionDeleteData  = "delete"
	bitmexActionUpdateData  = "update"

	WSTimeOut = 5 * time.Second
)

type SubscribeInfo struct {
	Op    string
	Param string
}

type BitmexWS struct {
	TableLen               int
	baseURL                string
	symbol                 string
	key                    string
	secret                 string
	proxy                  string
	wsConn                 *websocket.Conn
	partialLoadedTrades    bool
	partialLoadedOrderbook bool
	orderBook              OrderBookMap
	trades                 []*models.Trade
	pos                    PositionMap
	orders                 *OrderMap

	pongChan chan int
	shutdown *Shutdown

	lastDepth      Depth
	lastDepthMutex sync.RWMutex

	lastTrade      Trade
	lastTradeMutex sync.RWMutex

	lastPosition      []Position
	lastPositionMutex sync.RWMutex

	lastOrder      []Order
	lastOrderMutex sync.RWMutex

	tradeChan chan Trade
	depthChan chan Depth
	klineChan map[string]chan *Candle
	timer     *time.Timer

	subcribeTypes []SubscribeInfo
}

func NewBitmexWS(symbol, key, secret, proxy string) (bw *BitmexWS) {
	bw = NewBitmexWSWithURL(symbol, key, secret, proxy, bitmexWSURL)
	return
}

func NewBitmexWSTest(symbol, key, secret, proxy string) (bw *BitmexWS) {
	bw = NewBitmexWSWithURL(symbol, key, secret, proxy, testBitmexWSURL)
	return
}

func NewBitmexWSWithURL(symbol, key, secret, proxy, wsURL string) (bw *BitmexWS) {
	bw = new(BitmexWS)
	bw.baseURL = wsURL
	bw.symbol = symbol
	bw.key = key
	bw.secret = secret
	bw.proxy = proxy
	bw.orderBook = NewOrderBookMap()
	bw.pos = NewPositionMap()
	bw.orders = NewOrderMap()
	bw.pongChan = make(chan int, 1)
	bw.shutdown = NewRoutineManagement()
	bw.timer = time.NewTimer(WSTimeOut)
	bw.subcribeTypes = []SubscribeInfo{
		SubscribeInfo{Op: BitmexWSOrderbookL2, Param: bw.symbol},
		SubscribeInfo{Op: BitmexWSTrade, Param: bw.symbol},
		SubscribeInfo{Op: BitmexWSPosition, Param: bw.symbol},
		SubscribeInfo{Op: BitmexWSOrder, Param: bw.symbol},
	}
	bw.klineChan = make(map[string]chan *Candle)
	return
}

func (bw *BitmexWS) SetSymbol(symbol string) (err error) {
	bw.symbol = symbol
	return
}

func (bw *BitmexWS) SetKlineChan(binSize string, klineChan chan *Candle) (err error) {
	bw.klineChan[binSize] = klineChan
	return
}

func (bw *BitmexWS) SetSubscribe(subcribeTypes []SubscribeInfo) {
	bw.subcribeTypes = subcribeTypes
	return
}

func (bw *BitmexWS) AddSubscribe(subcribeInfo SubscribeInfo) {
	bw.subcribeTypes = append(bw.subcribeTypes, subcribeInfo)
}

func (bw *BitmexWS) SetProxy(proxy string) {
	bw.proxy = proxy
}

// SetLastDepth set depth data,call by websocket message handler
func (bw *BitmexWS) SetLastDepth(depth Depth) {
	bw.lastDepthMutex.Lock()
	bw.lastDepth = depth
	bw.lastDepthMutex.Unlock()
	return
}

// GetLastDepth get last depths
func (bw *BitmexWS) GetLastDepth() (depth Depth) {
	bw.lastDepthMutex.RLock()
	depth = bw.lastDepth
	bw.lastDepthMutex.RUnlock()
	return
}

// SetLastTrade set depth data,call by websocket message handler
func (bw *BitmexWS) SetLastTrade(trade Trade) {
	bw.lastTradeMutex.Lock()
	bw.lastTrade = trade
	bw.lastTradeMutex.Unlock()
	return
}

// GetLastDepth get last depths
func (bw *BitmexWS) GetLastTrade() (trade Trade) {
	bw.lastTradeMutex.RLock()
	trade = bw.lastTrade
	bw.lastTradeMutex.RUnlock()
	return
}

func (bw *BitmexWS) SetLastPos(pos []Position) {
	bw.lastPositionMutex.Lock()
	bw.lastPosition = pos
	bw.lastPositionMutex.Unlock()
}

func (bw *BitmexWS) GetLastPos() (poses []Position) {
	bw.lastPositionMutex.RLock()
	poses = bw.lastPosition
	bw.lastPositionMutex.RUnlock()
	// log.Debug("processPosition", poses)
	return
}

func (bw *BitmexWS) UpdateOrders(orders []Order) {
	tOrders := []*models.Order{}
	for _, order := range orders {
		tOrder := &models.Order{
			OrderID:   &order.OrderID,
			Currency:  order.Currency,
			OrderQty:  int64(order.Amount),
			AvgPx:     order.PriceAvg,
			Price:     order.Price,
			OrdStatus: order.Status,
			Side:      order.Side,
			Timestamp: strfmt.DateTime(order.Time),
		}
		tOrders = append(tOrders, tOrder)
	}
	bw.orders.Update(tOrders, false)
	bw.SetLastOrders(bw.orders.Orders())
}

func (bw *BitmexWS) SetLastOrders(orders []Order) {
	bw.lastOrderMutex.Lock()
	bw.lastOrder = orders
	bw.lastOrderMutex.Unlock()
}

func (bw *BitmexWS) GetLastOrders() (orders []Order) {
	bw.lastOrderMutex.RLock()
	orders = bw.lastOrder
	bw.lastOrderMutex.RUnlock()
	return
}

func (bw *BitmexWS) GetLastOrder(oid string) (order Order, err error) {
	flag := false
	bw.lastOrderMutex.RLock()
	for _, o := range bw.lastOrder {
		if o.OrderID == oid {
			order = o
			flag = true
			break
		}
	}
	bw.lastOrderMutex.RUnlock()
	if !flag {
		err = NoOrderFound
	}
	return
}

func (bw *BitmexWS) SetTradeChan(tradeChan chan Trade) {
	bw.tradeChan = tradeChan
}

func (bw *BitmexWS) SetDepthChan(depthChan chan Depth) {
	bw.depthChan = depthChan
}

func (bw *BitmexWS) Connect() (err error) {
	dialer := websocket.Dialer{}
	if bw.proxy != "" {
		uProxy, err := url.Parse(bw.proxy)
		if err != nil {
			return err
		}
		dialer.Proxy = http.ProxyURL(uProxy)
	}

	bw.wsConn, _, err = dialer.Dial(bw.baseURL, nil)
	if err != nil {
		return err
	}
	_, p, err := bw.wsConn.ReadMessage()
	if err != nil {
		return err
	}

	bw.partialLoadedOrderbook = false
	bw.partialLoadedTrades = false
	var welcome Welcome
	err = json.Unmarshal(p, &welcome)
	if err != nil {
		return err
	}
	log.Debug("welcome:", string(p))
	go bw.connectionHandler()
	go bw.handleMessage()
	err = bw.sendAuth()
	if err != nil {
		return err
	}
	err = bw.subscribe()
	if err != nil {
		return err
	}

	if bw.key != "" {
		err = bw.sendAuth()
		if err != nil {
			return
		}
	}
	return
}

// Timer handles connection loss or failure
func (bw *BitmexWS) connectionHandler() {
	defer func() {
		log.Debug("Bitmex websocket: Connection handler routine shutdown")
	}()

	shutdown := bw.shutdown.addRoutine()
	bw.timer.Reset(WSTimeOut)
	for {
		select {
		case <-bw.timer.C:
			timeout := time.After(WSTimeOut)
			log.Println("time out first,send ping...")
			err := bw.wsConn.WriteJSON("ping")
			if err != nil {
				bw.reconnect()
				return
			}
		OUT:
			for {
				select {
				case <-bw.pongChan:
					log.Debug("Bitmex websocket: PONG chan received")
					// bw.timer.Reset(WSTimeOut)
					time.Sleep(time.Microsecond)
					break OUT
				case <-timeout:
					log.Println("Bitmex websocket: Connection timed out - Closing connection....")
					bw.wsConn.Close()

					log.Println("Bitmex websocket: Connection timed out - Reconnecting...")
					bw.reconnect()
					return
				}
			}
		case <-shutdown:
			log.Println("Bitmex websocket: shutdown requested - Closing connection....")
			bw.wsConn.Close()
			log.Println("Bitmex websocket: Sending shutdown message")
			bw.shutdown.routineShutdown()
			return
		}
	}
}

// Reconnect handles reconnections to websocket API
func (bw *BitmexWS) reconnect() {
	for {
		err := bw.Connect()
		if err != nil {
			log.Println("Bitmex websocket: Connection timed out - Failed to connect, sleeping...")
			time.Sleep(time.Second * 2)
			continue
		}
		return
	}
}

// sendAuth sends an authenticated subscription
func (bw *BitmexWS) sendAuth() error {
	timestamp := time.Now().Add(time.Hour * 1).Unix()
	newTimestamp := strconv.FormatInt(timestamp, 10)
	hmac := GetHMAC(HashSHA256,
		[]byte("GET/realtime"+newTimestamp),
		[]byte(bw.secret))

	signature := hex.EncodeToString(hmac)

	var sendAuth WSCmd
	sendAuth.Command = "authKeyExpires"
	sendAuth.Args = append(sendAuth.Args, bw.key)
	sendAuth.Args = append(sendAuth.Args, timestamp)
	sendAuth.Args = append(sendAuth.Args, signature)
	return bw.wsConn.WriteJSON(sendAuth)
}

// subscribe subscribes to a websocket channel
func (bw *BitmexWS) subscribe() (err error) {
	// Subscriber
	var subscriber WSCmd
	subscriber.Command = "subscribe"

	// Announcement subscribe
	// subscriber.Args = append(subscriber.Args, bitmexWSAnnouncement)
	for _, v := range bw.subcribeTypes {
		subscriber.Args = append(subscriber.Args,
			v.Op+":"+v.Param)
	}
	err = bw.wsConn.WriteJSON(subscriber)
	if err != nil {
		return err
	}

	return nil
}

func (bw *BitmexWS) handleMessage() {
	var data []byte
	var msg string
	var err error
	for {
		_, data, err = bw.wsConn.ReadMessage()
		if err != nil {
			log.Error("Bitmex websocket read error:", err.Error())
			return
		}
		msg = string(data)
		if strings.Contains(msg, "pong") {
			log.Debug("Bitmex websocket pong received")
			bw.pongChan <- 1
			continue
		}
		if strings.Contains(msg, "ping") {
			err = bw.wsConn.WriteJSON("pong")
			if err != nil {
				log.Error("Bitmex websocket error:", err.Error())
			}
			continue
		}
		bw.timer.Reset(WSTimeOut)
		var ret Resp
		err = ret.Decode(data)
		if err != nil {
			log.Fatal(err)
			continue
		}
		if ret.HasStatus() {
			log.Error("Bitmex websocket error:", msg)
		} else if ret.HasSuccess() {
			if !ret.Success {
				log.Error("Bitmex websocket error:", msg)
			} else {
				log.Debug("Bitmex websocket subscribed success")
			}
		} else if ret.HasTable() {
			switch ret.Table {
			case BitmexWSOrderbookL2:
				err = bw.processOrderbook(&ret)
			case BitmexWSTrade:
				err = bw.processTrade(&ret)
			case BitmexWSAnnouncement:
				// err = bw.processTrade(&ret)
			case BitmexWSPosition:
				// log.Debug("processPosition", msg)
				err = bw.processPosition(&ret)
			case BitmexWSOrder:
				err = bw.processOrder(&ret)
			case BitmexWSTradeBin1m, BitmexWSTradeBin5m, BitmexWSTradeBin1h, BitmexWSTradeBin1d:
				err = bw.processTradeBin(ret.Table, &ret)
			default:
				log.Println(ret.Table, msg)
			}
		} else {

		}

	}
}

func (bw *BitmexWS) processOrderbook(msg *Resp) (err error) {
	datas := msg.GetOrderbookL2()
	updated := len(datas)
	switch msg.Action {
	case bitmexActionInitialData:
		if !bw.partialLoadedOrderbook {
			datas.GetDataToMap(bw.orderBook)
		}
		bw.partialLoadedOrderbook = true
		updated = 0
	case bitmexActionUpdateData:
		if bw.partialLoadedOrderbook {
			var v *OrderBookL2
			var ok bool
			for _, elem := range datas {
				v, ok = bw.orderBook[elem.Key()]
				if ok {
					// price is same while id is same
					// v.Price = elem.Price
					v.Side = elem.Side
					v.Size = elem.Size
					updated--
				}
			}

		}
	case bitmexActionInsertData:
		if bw.partialLoadedOrderbook {
			for _, elem := range datas {
				bw.orderBook[elem.Key()] = elem
				updated--
			}

		}
	case bitmexActionDeleteData:
		if bw.partialLoadedOrderbook {
			for _, elem := range datas {
				delete(bw.orderBook, elem.Key())
				updated--
			}

		}
	default:
		err = fmt.Errorf("unsupport action:%s", msg.Action)
		return
	}
	if updated != 0 {
		return errors.New("Bitmex websocket error: Elements not updated correctly")
	}
	depth := bw.orderBook.GetDepth()
	// log.Debug("depth:", depth)
	bw.SetLastDepth(depth)
	if bw.depthChan != nil {
		bw.depthChan <- depth
	}
	return
}

func (bw *BitmexWS) processTrade(msg *Resp) (err error) {
	var datas []*models.Trade
	switch msg.Action {
	case bitmexActionInitialData:
		if !bw.partialLoadedTrades {
			datas = msg.GetTradeData()
		}
		bw.partialLoadedTrades = true
	case bitmexActionInsertData:
		if bw.partialLoadedTrades {
			datas = append(bw.trades, msg.GetTradeData()...)
		}
	default:
		err = fmt.Errorf("Bitmex websocket error: unsupport action:%s", msg.Action)
	}
	if err != nil {
		return
	}
	if len(datas) > bw.TableLen {
		bw.trades = datas[len(datas)-bw.TableLen:]
	} else {
		bw.trades = datas
	}
	if len(datas) > 0 {
		v := datas[len(datas)-1]
		bw.SetLastTrade(transTrade(v))
	}
	if bw.tradeChan != nil {
		for _, v := range datas {
			bw.tradeChan <- transTrade(v)
		}
	}
	return
}

func (bw *BitmexWS) processPosition(msg *Resp) (err error) {
	datas := msg.GetPostions()
	switch msg.Action {
	case bitmexActionInitialData, bitmexActionUpdateData, bitmexActionInsertData:
		bw.pos.Update(datas)
	// case bitmexActionDeleteData:
	default:
		err = fmt.Errorf("unsupport action:%s", msg.Action)
		return
	}
	bw.SetLastPos(bw.pos.Pos())
	return
}

func (bw *BitmexWS) processOrder(msg *Resp) (err error) {
	datas := msg.GetOrders()
	switch msg.Action {
	case bitmexActionInitialData, bitmexActionUpdateData, bitmexActionInsertData:
		bw.orders.Update(datas, false)
	case bitmexActionDeleteData:
		bw.orders.Update(datas, true)
	default:
		err = fmt.Errorf("unsupport action:%s", msg.Action)
		return
	}
	bw.SetLastOrders(bw.orders.Orders())
	return
}

func (bw *BitmexWS) processTradeBin(binSize string, msg *Resp) (err error) {
	klineChan, ok := bw.klineChan[binSize]
	if !ok {
		log.Debug("no such kline chan", binSize)
		return
	}
	datas := msg.GetTradeBin()
	switch msg.Action {
	case bitmexActionInitialData, bitmexActionUpdateData, bitmexActionInsertData:
		var candles []*Candle
		transCandle(datas, &candles, binSize)
		for _, v := range candles {
			klineChan <- v
		}
	// case bitmexActionDeleteData:
	default:
		err = fmt.Errorf("unsupport action:%s", msg.Action)
		return
	}
	return
}
