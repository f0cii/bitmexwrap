package coinex

import (
	"time"

	. "github.com/sumorf/trademodel"
)

const (
	Long  = 1
	Short = 2
)

type Contract struct {
	Symbol string
	Name   string
	Expiry time.Time
}

func (c Contract) String() string {
	return c.Symbol + " " + c.Name
}

type Position struct {
	Info        Contract // 合约信息
	Type        int      // 合约类型，Long: 多头，Short: 空头
	Hold        float64  // 持有仓位
	Price       float64  // 开仓价格/开仓均价 AvgEntryPrice
	CostPrice   float64  // 成本价格
	LastPrice   float64  // 最新成交价
	MarkPrice   float64  // 标记价格
	ProfitRatio float64  // 盈利比例,正数表示盈利，负数表示亏损
}

type BaseExchanger interface {
	Buy(price float64, amount float64) (*Order, error)
	Sell(price float64, amount float64) (*Order, error)
}

type FuturesBaseExchanger interface {
	BaseExchanger
	OpenLong(price float64, amount float64) (*Order, error)
	CloseLong(price float64, amount float64) (*Order, error)
	OpenShort(price float64, amount float64) (*Order, error)
	CloseShort(price float64, amount float64) (*Order, error)
	OpenLongMarket(amount float64) (*Order, error)
	CloseLongMarket(amount float64) (*Order, error)
	OpenShortMarket(amount float64) (*Order, error)
	CloseShortMarket(amount float64) (*Order, error)
}

type FuturesExchanger interface {
	FuturesBaseExchanger

	Contracts() ([]Contract, error)
	Positions() ([]Position, error)
	ContractBalances() (map[Contract]Balance, error)

	Depth(int) (Orderbook, error)
	Ticker() (Ticker, error)

	SetSymbol(symbol string) error
	SetContract(contract string) error
	SetLever(lever float64) error

	KlineRecent(count int32, binSize string) (klines []*Candle, err error)
	Kline(start, end time.Time, nLimit int, binSize string) (klines []*Candle, err error)
	KlineChan(start, end time.Time, bSize string) (klines chan []interface{}, err error)
}
