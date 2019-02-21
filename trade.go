package bitmexwrap

import "time"

type ActionType int

const (
	ActionBuy ActionType = iota
	ActionSell
	ActionOpenLong
	ActionCloseLong
	ActionOpenShort
	ActionCloseShort
)

func (at ActionType) String() (msg string) {
	switch at {
	case ActionBuy:
		msg = "buy"
	case ActionSell:
		msg = "sell"
	case ActionOpenLong:
		msg = "openLong"
	case ActionCloseLong:
		msg = "closeLong"
	case ActionOpenShort:
		msg = "openShort"
	case ActionCloseShort:
		msg = "closeShort"
	default:
	}
	return
}

type Trade struct {
	ID     string
	Symbol string
	Time   time.Time
	Price  float64
	Amount float64
	Side   string
	Remark string
}

type TradeAction struct {
	Action ActionType
	Amount float64
	Price  float64
	Time   time.Time
}

func (ta *TradeAction) IsBuy() bool {
	switch ta.Action {
	case ActionBuy, ActionOpenLong, ActionCloseShort:
		return true
	case ActionSell, ActionCloseLong, ActionOpenShort:
		return false
	default:
	}
	return false
}
