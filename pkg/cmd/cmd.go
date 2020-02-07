package cmd

import (
	"fmt"
	"github.com/goecology/muses/pkg/common"
	"os"
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
	c.initStart()
	return nil
}

func (c *callerStore) InitCaller() error {
	if err := common.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}
