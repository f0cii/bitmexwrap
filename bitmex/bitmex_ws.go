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

	. "github.com/SuperGod/coinex"

	. "github.com/SuperGod/trademodel"

	"github.com/SuperGod/coinex/bitmex/models"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	MaxTableLen     = 200
	bitmexWSURL     = "wss://www.bitmex.com/realtime"
	testBitmexWSURL = "wss://testnet.bitmex.com/realtime"

	bitmexWSOrderbookL2  = "orderBookL2"
	bitmexWSOrderbookL10 = "orderBook10"
	bitmexWSTrade        = "trade"
	bitmexWSAnnouncement = "announcement"
	bitmexWSLiquidation  = "liquidation"

	bitmexWSOrder    = "order"
	bitmexWSPosition = "position"

	bitmexActionInitialData = "partial"
	bitmexActionInsertData  = "insert"
	bitmexActionDeleteData  = "delete"
	bitmexActionUpdateData  = "update"

	WSTimeOut = 5 * time.Second
)

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

	pongChan chan int
	shutdown *Shutdown

	lastDepth      Depth
	lastDepthMutex sync.RWMutex

	lastTrade      Trade
	lastTradeMutex sync.RWMutex

	lastPosition      []Position
	lastPositionMutex sync.RWMutex

	tradeChan chan Trade
	depthChan chan Depth
	timer     *time.Timer
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
	bw.pongChan = make(chan int, 1)
	bw.shutdown = NewRoutineManagement()
	bw.timer = time.NewTimer(WSTimeOut)
	return
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

// SetLastDepth set depth data,call by websocket message handler
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
	subscriber.Args = append(subscriber.Args, bitmexWSAnnouncement)
	subscriber.Args = append(subscriber.Args,
		bitmexWSOrderbookL2+":"+bw.symbol)
	subscriber.Args = append(subscriber.Args,
		bitmexWSTrade+":"+bw.symbol)
	subscriber.Args = append(subscriber.Args,
		bitmexWSPosition+":"+bw.symbol)

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
			case bitmexWSOrderbookL2:
				err = bw.processOrderbook(&ret)
			case bitmexWSTrade:
				err = bw.processTrade(&ret)
			case bitmexWSAnnouncement:
				// err = bw.processTrade(&ret)
			case bitmexWSPosition:
				// log.Debug("processPosition", msg)
				err = bw.processPosition(&ret)
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
					v.Price = elem.Price
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

	if bw.depthChan != nil {
		depth := bw.orderBook.GetDepth()
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
