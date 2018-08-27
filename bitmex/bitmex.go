package bitmex

import (
	"context"
	"fmt"
	"net/url"
	"time"

	. "github.com/SuperGod/coinex"
	apiclient "github.com/SuperGod/coinex/bitmex/client"
	"github.com/SuperGod/coinex/bitmex/client/instrument"
	"github.com/SuperGod/coinex/bitmex/client/order_book"
	"github.com/SuperGod/coinex/bitmex/client/position"
	"github.com/SuperGod/coinex/bitmex/client/trade"
	apiuser "github.com/SuperGod/coinex/bitmex/client/user"
	"github.com/SuperGod/coinex/bitmex/models"
	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"
)

const (
	BaseURL     = "www.bitmex.com"
	TestBaseURL = "testnet.bitmex.com"
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

func (b *Bitmex) context() (ctx context.Context) {
	ctx = context.Background()
	return
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
	pos, err := b.api.Position.PositionGet(&position.PositionGetParams{}, nil)
	if err != nil {
		return
	}
	var orderType int
	for _, v := range pos.Payload {
		if v.CurrentQty > 0 {
			orderType = Long
		} else {
			orderType = Short
		}
		if v.CurrentQty == 0 {
			continue
		}
		// UnrealisedRoePcnt 是按标记价格计算的盈亏
		positions = append(positions, Position{Info: Contract{Symbol: *v.Symbol, Name: *v.Symbol},
			Type:        orderType,
			Hold:        float64(v.CurrentQty),
			ProfitRatio: float64(v.UnrealisedRoePcnt),
		},
		)
	}
	return
}

// ContractBalances get balances of each contract
func (b *Bitmex) ContractBalances() (balances map[Contract]Balance, err error) {
	wallet, err := b.api.User.UserGetWallet(&apiuser.UserGetWalletParams{}, nil)
	if err != nil {
		return
	}
	fmt.Println(wallet)
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
func (b *Bitmex) Depth(d int) (depth Depth, err error) {
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
func (b *Bitmex) GetDepth(d int) (depth Depth, err error) {
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
	ticker.Ask = depth.Buys[0].Price
	ticker.Bid = depth.Sells[0].Price
	return
}

// Ticker
func (b *Bitmex) GetTicker() (ticker Ticker, err error) {
	// symbol := optional.NewString(b.symbol)
	reverse := true
	nCount := int32(10)
	// oCount := optional.NewFloat32(float32(nCount))
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
	ticker.Volume = v.HomeNotional
	ticker.CurrencyPair = b.symbol
	ticker.Ask = depth.Sells[0].Price
	ticker.Bid = depth.Buys[0].Price
	return
}

// SetSymbol set symbol
func (b *Bitmex) SetSymbol(symbol string) (err error) {
	b.symbol = symbol
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

func (b *Bitmex) Kline(start, end time.Time, bSize string) (klines []*models.TradeBin, err error) {
	startTime := strfmt.DateTime(start)
	endTime := strfmt.DateTime(end)
	var nStart int32
	var nCount int32
	var nRet int32
	params := &trade.TradeGetBucketedParams{BinSize: &bSize, StartTime: &startTime, EndTime: &endTime}
	for {
		params.Start = &nStart
		params.Count = &nCount
		klineInfo, err := b.api.Trade.TradeGetBucketed(params)
		if err != nil {
			break
		}
		klines = append(klines, klineInfo.Payload...)
		nRet = int32(len(klineInfo.Payload))
		nStart += nRet
		if nRet < nCount {
			break
		}
	}
	return
}

// KlineRecent get recent nCount klines
func (b *Bitmex) KlineRecent(nCount int32, bSize string) (klines []*models.TradeBin, err error) {
	bReverse := true
	params := &trade.TradeGetBucketedParams{BinSize: &bSize, Count: &nCount, Reverse: &bReverse}
	klineInfo, err := b.api.Trade.TradeGetBucketed(params)
	if err != nil {
		return
	}
	klines = klineInfo.Payload
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
