package cmd

import (
	"github.com/goecology/muses/pkg/common"
)

var defaultCaller = &callerStore{
	Name:        common.ModCmdName,
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
	//c.initStart()
	return nil
}

func (c *callerStore) InitCaller() error {

	return nil
}
