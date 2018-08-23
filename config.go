package coinex

import (
	"encoding/json"
	"os"
	"os/user"
)

type ConfigItem struct {
	Type   string
	Key    string
	Secret string
}

type Config struct {
	Exchanges map[string]ConfigItem
	Proxy     string
}

func LoadConfigs() (configs Config, err error) {
	u, err := user.Current()
	if err != nil {
		return
	}
	return LoadConfigFile(u.HomeDir + "/.config/exchange.json")
}

func LoadConfigFile(cfg string) (configs Config, err error) {
	f, err := os.Open(cfg)
	if err != nil {
		return
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	err = dec.Decode(&configs)
	return
}

func (c *Config) Get(name string) (key, secret string) {
	ex, ok := c.Exchanges[name]
	if !ok {
		return
	}
	key = ex.Key
	secret = ex.Secret
	return
}
