package cmd

import (
	"github.com/goecology/muses/pkg/common"
	"github.com/goecology/muses/pkg/system"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: system.BuildInfo.Name,
}

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
	return nil
}

func (c *callerStore) InitCaller() error {
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}
