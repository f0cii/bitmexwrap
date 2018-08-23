package coinex

import "time"

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
	ProfitRatio float64  // 盈利比例,正数表示盈利，负数表示亏岁
}

type FuturesExchange interface {
	Contracts() ([]Contract, error)
	Positions() ([]Position, error)
	ContractBalances() (map[Contract]Balance, error)

	Depth(int) (Depth, error)
	Ticker() (Ticker, error)

	SetSymbol(symbol string) error
	SetContract(contract string) error
	SetLever(lever int) error

	OpenLong(price float64, amount float64) error
	CloseLong(price float64, amount float64) error
	OpenShort(price float64, amount float64) error
	CloseShort(price float64, amount float64) error
}
