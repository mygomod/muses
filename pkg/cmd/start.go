package cmd

import (
	"github.com/goecology/muses/pkg/common"
	"github.com/goecology/muses/pkg/system"
	"github.com/spf13/cobra"
)

var startCmd *cobra.Command

func InitCommand(f func(cmd *cobra.Command, args []string) error) {
	startCmd = &cobra.Command{
		Use:  "start",
		Long: `Starts ` + system.BuildInfo.Name + ` server`,
		RunE: f,
	}
	startCmd.PersistentFlags().StringVarP(&common.CmdConfigPath, "conf", "c", "conf/conf.toml", "conf path")
	common.RootCmd.AddCommand(startCmd)
}
