package cmd

import (
	"github.com/goecology/muses/pkg/common"
	"github.com/goecology/muses/pkg/system"
	"github.com/spf13/cobra"
)

var startCmd *cobra.Command
var startFn func(cmd *cobra.Command, args []string)

func (c *callerStore) initStart() {
	startCmd = &cobra.Command{
		Use:  "start",
		Long: `Starts ` + system.BuildInfo.Name + ` server`,
		Run:  startFn,
	}
	startCmd.PersistentFlags().StringVarP(&common.CmdConfigPath, "conf", "c", "conf/conf.toml", "conf path")
	common.RootCmd.AddCommand(startCmd)
}

func SetStartFn(f func(cmd *cobra.Command, args []string)) {
	startFn = f
}
