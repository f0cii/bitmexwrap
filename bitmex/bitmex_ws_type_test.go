package bitmex

import (
	"testing"
)

func TestA(t *testing.T) {
	data := `{"table":"orderBook10","action":"update","data":[{"symbol":"XBTUSD","bids":[[3891.5,1223796],[3891,304394],[3890.5,41041],[3890,281329],[3889.5,303518],[3889,334669],[3888.5,325687],[3888,330918],[3887.5,51826],[3887,131232]],"timestamp":"2019-02-20T07:31:13.093Z","asks":[[3892,238358],[3892.5,164616],[3893,79037],[3893.5,61063],[3894,201077],[3894.5,215530],[3895,562813],[3895.5,222641],[3896,339422],[3896.5,262947]]}]}`
	var ret Resp
	err := ret.Decode([]byte(data))
	if err != nil {
		t.Error(err)
		return
	}
	orderbooks := ret.GetOrderbook10()
	t.Log(orderbooks)
}
