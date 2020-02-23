package app

import (
	"github.com/BurntSushi/toml"
	"github.com/i2eco/muses/pkg/common"
)

var defaultCaller = &callerStore{
	Name:        common.ModAppName,
	IsNecessary: true,
}

type callerStore struct {
	Name         string
	IsNecessary  bool
	IsBackground bool
	caller       *Client
	cfg          Cfg
}

type Client struct{}

func Register() common.Caller {
	return defaultCaller
}

func Config() Cfg {
	return defaultCaller.cfg
}

func (c *callerStore) InitCfg(cfg []byte) error {
	if err := toml.Unmarshal(cfg, &c.cfg); err != nil {
		return err
	}
	return nil
}

func (c *callerStore) InitCaller() error {
	return nil
}
