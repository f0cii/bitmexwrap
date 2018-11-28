package bitmexwrap

import (
	"fmt"
	"time"
)

type Candle struct {
	ID     int64   `xorm:"pk autoincr null 'id'"`
	Start  int64   `xorm:"unique index 'start'"`
	Open   float64 `xorm:"notnull 'open'"`
	High   float64 `xorm:"notnull 'high'"`
	Low    float64 `xorm:"notnull 'low'"`
	Close  float64 `xorm:"notnull 'close'"`
	VWP    float64 `xorm:"notnull 'vwp'"`
	Volume float64 `xorm:"notnull 'volume'"`
	Trades int64   `xorm:"notnull 'trades'"`
	Table  string  `xorm:"-"`
}

func (c Candle) TableName() string {
	return c.Table
}

func (c Candle) Time() time.Time {
	return time.Unix(c.Start, 0)
}

func (c Candle) String() string {
	return fmt.Sprintf("%s open:%f close:%f low:%f high:%f volume:%f vwp:%f trades:%d", c.Time().String(), c.Open, c.Close, c.Low, c.High, c.Volume, c.VWP, c.Trades)
}
