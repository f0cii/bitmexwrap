package bitmex

import (
	"sort"
	"testing"
	"time"
)

func TestWSOrderBook(t *testing.T) {
	b := GetClient()
	b.maxLocalDepth = 0
	err := b.StartWS()
	if err != nil {
		t.Fatal(err.Error())
	}
	time.Sleep(10 * time.Second)

	depth, err := b.Depth(0)
	if err != nil {
		t.Fatal(err.Error())
	}
	// t.Log("depth time:", depth.UpdateTime)
	// t.Log("==========buy======")
	// for _, v := range depth.Buys {
	// 	t.Log(v.Price, v.Amount)
	// }
	// t.Log("==========sell======")
	// for _, v := range depth.Sells {
	// 	t.Log(v.Price, v.Amount)
	// }

	sort.Slice(depth.Sells, func(i, j int) bool {
		return depth.Sells[i].Amount > depth.Sells[j].Amount
	})
	for k, v := range depth.Sells[0:10] {
		t.Logf("sort %d %f %f\n", k, v.Price, v.Amount)
	}
	t.Logf("depth:%d %d\n", len(depth.Buys), len(depth.Sells))
	t.Logf("depth price:%f %f\n", depth.Buys[0].Price, depth.Sells[0].Price)
}

func TestWSTrades(t *testing.T) {
	b := GetClient()
	err := b.StartWS()
	if err != nil {
		t.Fatal(err.Error())
	}
	time.Sleep(5 * time.Second)
	depth, err := b.Depth(0)
	if err != nil {
		t.Fatal(err.Error())
	}
	// t.Log("depth time:", depth.UpdateTime)
	// t.Log("==========buy======")
	// for _, v := range depth.Buys {
	// 	t.Log(v.Price, v.Amount)
	// }
	// t.Log("==========sell======")
	// for _, v := range depth.Sells {
	// 	t.Log(v.Price, v.Amount)
	// }
	t.Logf("depth:%d %d\n", len(depth.Buys), len(depth.Sells))
	t.Logf("depth price:%f %f\n", depth.Buys[0].Price, depth.Sells[0].Price)
}
