package bitmex

import (
	"fmt"
	"net/url"
	"time"

	. "github.com/sumorf/coinex"
	. "github.com/sumorf/trademodel"

	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
	apiclient "github.com/sumorf/coinex/bitmex/client"
	"github.com/sumorf/coinex/bitmex/client/instrument"
	"github.com/sumorf/coinex/bitmex/client/order_book"
	"github.com/sumorf/coinex/bitmex/client/position"
	"github.com/sumorf/coinex/bitmex/client/trade"
	apiuser "github.com/sumorf/coinex/bitmex/client/user"
	"github.com/sumorf/coinex/bitmex/models"
)

const (
	BaseURL     = "www.bitmex.com"
	TestBaseURL = "testnet.bitmex.com"
)

var (
	timeLocal, _ = time.LoadLocation("Local")
)

type Bitmex struct {
	wsAPI     *BitmexWS
	api       *apiclient.APIClient
	symbol    string
	lever     float64
	APIKey    string
	APISecret string
	proxy     string
	enableWS  bool
	trans     *Transport
}

func NewBitmex(key, secret string) (b *Bitmex) {
	// return NewBitmexFromCfg(key, secret, BaseURL, bitmexWSURL, nil)
	b = NewBitmexFromCfg(key, secret, BaseURL, nil)
	b.wsAPI = NewBitmexWS(b.symbol, key, secret, "")
	return
}

func NewBitmexTest(key, secret string) (b *Bitmex) {
	// return NewBitmexFromCfg(key, secret, TestBaseURL, testBitmexWSURL, nil)
	b = NewBitmexFromCfg(key, secret, TestBaseURL, nil)
	b.wsAPI = NewBitmexWSTest(b.symbol, key, secret, "")
	return
}

// func NewBitmexFromCfg(key, secret, baseURL string, cfg *apiclient.Configuration) *Bitmex {
func NewBitmexFromCfg(key, secret, baseURL string, cfg *apiclient.TransportConfig) *Bitmex {

	b := new(Bitmex)
	if cfg == nil {
		cfg = &apiclient.TransportConfig{}
	}
	cfg.Host = baseURL
	cfg.BasePath = "/api/v1"
	cfg.Schemes = []string{"https"}
	b.api = apiclient.NewHTTPClientWithConfig(nil, cfg)
	b.trans = NewTransport(cfg.Host, cfg.BasePath, key, secret, cfg.Schemes)
	b.api.SetTransport(b.trans)
	b.symbol = "XBTUSD"
	b.APIKey = key
	b.APISecret = secret
	return b
}

func (b *Bitmex) SetDebug(bDebug bool) {
	b.trans.SetDebug(bDebug)
}

// SetMaxLocalDepth set max local depth cache len
func (b *Bitmex) SetMaxLocalDepth(nMaxDepth int) {
	b.wsAPI.TableLen = nMaxDepth
}

// StartWS start websocket connection
func (b *Bitmex) StartWS() (err error) {
	b.enableWS = true
	err = b.wsAPI.Connect()
	if err != nil {
		return
	}
	return
}

// WS return ws instance
func (b *Bitmex) WS() *BitmexWS {
	return b.wsAPI
}

func (b *Bitmex) SetTradeChan(tradeChan chan Trade) {
	b.wsAPI.SetTradeChan(tradeChan)
}

func (b *Bitmex) SetDepthChan(depthChan chan Depth) {
	b.wsAPI.SetDepthChan(depthChan)
}

// SetProxy set proxy of websocket
// example: socks5://127.0.0.1:1080
//          http://127.0.0.1:1080
func (b *Bitmex) SetProxy(proxy string) (err error) {
	_, err = url.Parse(proxy)
	if err != nil {
		return
	}
	b.wsAPI.SetProxy(proxy)
	return
}

// Contracts get all support contracts
// Fixme: "parse error", may be the swagger code wrong
func (b *Bitmex) Contracts() (contracts []Contract, err error) {
	ret, err := b.api.Instrument.InstrumentGetActiveAndIndices(&instrument.InstrumentGetActiveAndIndicesParams{})
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	for _, v := range ret.Payload {
		if v.State != "Open" {
			continue
		}
		contracts = append(contracts,
			Contract{Symbol: v.RootSymbol,
				Name:   *v.Symbol,
				Expiry: time.Time(v.Expiry)})
	}
	return
}

// Positions get current positions
func (b *Bitmex) Positions() (positions []Position, err error) {
	if b.enableWS {
		positions = b.wsAPI.GetLastPos()
		return
	}
	return b.GetPositions()
}

// GetPositions get current positions
func (b *Bitmex) GetPositions() (positions []Position, err error) {
	var pos *position.PositionGetOK
	pos, err = b.api.Position.PositionGet(&position.PositionGetParams{}, nil)
	if err != nil {
		return
	}
	var position *Position
	for _, v := range pos.Payload {
		log.Printf("%#v", v)
		position = transPosition(v)
		if position == nil {
			continue
		}
		// UnrealisedRoePcnt 是按标记价格计算的盈亏
		// UnrealisedPnl 未实现盈亏
		// UnrealisedPnlPcnt 未实现盈亏%
		// markPrice 标记价
		// avgEntryPrice 开仓均价
		positions = append(positions, *position)
	}
	return
}

// UpdatePositions update current positions
func (b *Bitmex) UpdatePositions() (err error) {
	var pos *position.PositionGetOK
	pos, err = b.api.Position.PositionGet(&position.PositionGetParams{}, nil)
	if err != nil {
		return
	}
	positions := pos.Payload
	b.wsAPI.pos.Update(positions)
	b.wsAPI.SetLastPos(b.wsAPI.pos.Pos())
	return
}

// ContractBalances get balances of each contract
func (b *Bitmex) ContractBalances() (balances map[Contract]Balance, err error) {
	wallet, err := b.api.User.UserGetWallet(&apiuser.UserGetWalletParams{}, nil)
	if err != nil {
		return
	}
	// fmt.Println(wallet)
	c := Contract{}
	balances = make(map[Contract]Balance)
	balance := Balance{
		Available: float64(wallet.Payload.Amount),
		Balance:   float64(wallet.Payload.Amount),
		Currency:  *wallet.Payload.Currency,
		Frozen:    0,
	}
	balances[c] = balance
	return
}

// Balance get balance of user
func (b *Bitmex) Balance() (balance Balance, err error) {
	var wallet *apiuser.UserGetWalletOK
	wallet, err = b.api.User.UserGetWallet(&apiuser.UserGetWalletParams{}, nil)
	if err != nil {
		return
	}
	balance.Available = float64(wallet.Payload.Amount)
	balance.Balance = float64(wallet.Payload.Amount)
	balance.Currency = *wallet.Payload.Currency
	balance.Frozen = 0
	return
}

// GetWalletHistory get wallet history
func (b *Bitmex) GetWalletHistory() (trans []*models.Transaction, err error) {
	var walletHistory *apiuser.UserGetWalletHistoryOK
	currency := "XBt"
	params := apiuser.UserGetWalletHistoryParams{
		Currency: &currency,
	}
	walletHistory, err = b.api.User.UserGetWalletHistory(&params, nil)
	if err != nil {
		return
	}
	trans = walletHistory.Payload
	return
}

func (b *Bitmex) User() (user *models.User, err error) {
	userInfo, err := b.api.User.UserGet(&apiuser.UserGetParams{}, nil)
	if err != nil {
		return
	}
	user = userInfo.Payload
	return
}

// Depth get depth
// if d is 0, get all depth
func (b *Bitmex) Depth(d int) (depth Orderbook, err error) {
	if !b.enableWS {
		return b.GetDepth(d)
	}
	lastDepth := b.wsAPI.GetLastDepth()
	if d > 0 && d < len(lastDepth.Buys) {
		depth.Buys = lastDepth.Buys[0:d]
	} else {
		depth.Buys = lastDepth.Buys
	}
	if d > 0 && d < len(lastDepth.Sells) {
		depth.Sells = lastDepth.Sells[0:d]
	} else {
		depth.Sells = lastDepth.Sells
	}
	depth.UpdateTime = lastDepth.UpdateTime
	return
}

// GetDepth get depth use RESTful API
func (b *Bitmex) GetDepth(d int) (depth Orderbook, err error) {
	nDepth := int32(d)
	ret, err := b.api.OrderBook.OrderBookGetL2(&order_book.OrderBookGetL2Params{Depth: &nDepth, Symbol: b.symbol})
	if err != nil {
		return
	}
	for _, v := range ret.Payload {
		if *v.Side == "Sell" {
			depth.Sells = append(depth.Sells,
				DepthInfo{Price: float64(v.Price),
					Amount: float64(v.Size)})
		} else {
			depth.Buys = append(depth.Buys,
				DepthInfo{Price: float64(v.Price),
					Amount: float64(v.Size)})
		}
	}
	depth.UpdateTime = time.Now()
	return
}

func (b *Bitmex) Price() (price float64, err error) {
	if !b.enableWS {
		var ticker Ticker
		ticker, err = b.GetTicker()
		if err != nil {
			return
		}
		price = ticker.Last
		return
	}
	trade := b.wsAPI.GetLastTrade()
	price = trade.Price
	return
}

func (b *Bitmex) Ticker() (ticker Ticker, err error) {
	if !b.enableWS {
		return b.GetTicker()
	}
	trade := b.wsAPI.GetLastTrade()
	depth := b.wsAPI.GetLastDepth()
	ticker.Last = trade.Price
	// ticker.Volume = trade.H
	if len(depth.Buys) > 0 {
		ticker.Ask = depth.Buys[0].Price
	}
	if len(depth.Sells) > 0 {
		ticker.Bid = depth.Sells[0].Price
	}
	return
}

// GetTicker get ticker
func (b *Bitmex) GetTicker() (ticker Ticker, err error) {
	reverse := true
	nCount := int32(10)
	ret2, err := b.api.Trade.TradeGet(&trade.TradeGetParams{Count: &nCount, Symbol: &b.symbol, Reverse: &reverse})
	if err != nil {
		return
	}
	if len(ret2.Payload) != int(nCount) {
		err = fmt.Errorf("trade result count not match:%d", len(ret2.Payload))
		return
	}
	depth, err := b.GetDepth(2)
	if err != nil {
		return
	}

	v := ret2.Payload[len(ret2.Payload)-1]
	ticker.Last = v.Price
	tm := time.Time(*v.Timestamp).In(timeLocal)
	ticker.Timestamp = tm
	ticker.Volume = v.HomeNotional
	ticker.CurrencyPair = b.symbol
	ticker.Ask = depth.Sells[0].Price
	ticker.Bid = depth.Buys[0].Price
	return
}

// SetSymbol set symbol
func (b *Bitmex) SetSymbol(symbol string) (err error) {
	b.symbol = symbol
	err = b.wsAPI.SetSymbol(symbol)
	return
}

func (b *Bitmex) SetContract(contract string) (err error) {
	return
}

// SetLever set contract lever
func (b *Bitmex) SetLever(lever float64) (err error) {
	b.lever = lever
	pos, err := b.api.Position.PositionUpdateLeverage(&position.PositionUpdateLeverageParams{Symbol: b.symbol, Leverage: lever}, nil)
	if err != nil {
		return
	}
	log.Println("set lever:", pos)
	return
}

func (b *Bitmex) GetLever() (lever float64, err error) {
	lever = b.lever
	return
}

// Kline Timestamp of kline is the end of the binSize
func (b *Bitmex) Kline(start, end time.Time, nLimit int, bSize string) (klines []*Candle, err error) {
	startTime := strfmt.DateTime(start)
	endTime := strfmt.DateTime(end)
	var nStart int32
	var nCount int32
	// var nRet int32
	nCount = 500
	if int32(nLimit) < nCount {
		nCount = int32(nLimit)
	}
	params := &trade.TradeGetBucketedParams{Symbol: &b.symbol, BinSize: &bSize, StartTime: &startTime, EndTime: &endTime}
	// for {
	params.Start = &nStart
	params.Count = &nCount
	klineInfo, err := b.api.Trade.TradeGetBucketed(params)
	if err != nil {
		// break
		return
	}
	transCandle(klineInfo.Payload, &klines, bSize)

	// nRet = int32(len(klineInfo.Payload))
	// nStart += nRet
	// if nRet < nCount {
	// 	break
	// }
	// }
	return
}

// KlineRecent get recent nCount klines
func (b *Bitmex) KlineRecent(nCount int32, bSize string) (klines []*Candle, err error) {
	bReverse := true
	params := &trade.TradeGetBucketedParams{BinSize: &bSize, Count: &nCount, Reverse: &bReverse, Symbol: &b.symbol}
	klineInfo, err := b.api.Trade.TradeGetBucketed(params)
	if err != nil {
		return
	}
	transCandle(klineInfo.Payload, &klines, bSize)
	return
}

// Trades get trades
func (b *Bitmex) Trades(start, end time.Time) (trades []Trade, err error) {
	startTime := strfmt.DateTime(start)
	endTime := strfmt.DateTime(end)
	var nStart, nCount, nRet int32
	nCount = 500
	params := &trade.TradeGetParams{StartTime: &startTime, EndTime: &endTime}
	for {
		params.Start = &nStart
		tradeInfo, err := b.api.Trade.TradeGet(params)
		if err != nil {
			break
		}
		nRet = int32(len(tradeInfo.Payload))
		for _, v := range tradeInfo.Payload {
			trades = append(trades, transTrade(v))
		}
		nStart += nRet
		if nRet < nCount {
			break
		}
	}
	return
}

func (b *Bitmex) TradesChan(start, end time.Time) (trades chan []interface{}, err error) {
	paramFunc := func() DownParam {
		p := trade.NewTradeGetParams()
		return p
	}
	downFunc := func(param DownParam) (data []interface{}, err1 error) {
		tradeParams := param.(*trade.TradeGetParams)
		trades, err1 := b.api.Trade.TradeGet(tradeParams)
		if err1 != nil {
			return
		}
		for _, v := range trades.Payload {
			data = append(data, v)
		}
		return
	}
	d := NewDataDownload(strfmt.DateTime(start), strfmt.DateTime(end), paramFunc, downFunc, 500, 10)
	trades = d.Start()
	return
}

func (b *Bitmex) KlineChan(start, end time.Time, bSize string) (klines chan []interface{}, err error) {
	paramFunc := func() DownParam {
		params := &trade.TradeGetBucketedParams{BinSize: &bSize, Symbol: &b.symbol}
		return params
	}
	downFunc := func(param DownParam) (data []interface{}, err1 error) {
		params := param.(*trade.TradeGetBucketedParams)
		klineInfo, err1 := b.api.Trade.TradeGetBucketed(params)
		if err1 != nil {
			return
		}
		for _, v := range klineInfo.Payload {
			data = append(data, transOneCandle(v))
		}
		return
	}
	d := NewDataDownload(strfmt.DateTime(start), strfmt.DateTime(end), paramFunc, downFunc, 500, 10)
	klines = d.Start()
	return
}
