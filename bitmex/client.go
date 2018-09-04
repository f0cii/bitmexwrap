package bitmex

import (
	. "github.com/SuperGod/coinex"
)

func GetClientByName(name string, bTest bool) (bm *Bitmex) {
	configs, err := LoadConfigs()
	if err != nil {
		panic(err.Error())
	}
	key, secret := configs.Get(name)
	if bTest {
		bm = NewBitmexTest(key, secret)
		bm.SetDebug(true)
	} else {
		bm = NewBitmex(key, secret)
	}
	bm.SetProxy(configs.Proxy)
	return bm
}

func GetDefaultClient() *Bitmex {
	return GetClientByName("bitmex", false)
}
