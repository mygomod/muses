package cmd

import (
	"github.com/mygomod/muses/pkg/system"
	"github.com/spf13/cobra"
)

var ConfigPath string
var Addr string
var StatAddr string

var startCmd *cobra.Command

func InitStartCommand(f func(cmd *cobra.Command, args []string) error) *cobra.Command {
	if startCmd != nil {
		return startCmd
	}
	startCmd = &cobra.Command{
		Use:  "start",
		Long: `Starts ` + system.BuildInfo.Name + ` server`,
		RunE: f,
	}
	startCmd.PersistentFlags().StringVarP(&ConfigPath, "conf", "c", "conf/conf.toml", "conf path")
	startCmd.PersistentFlags().StringVar(&Addr, "addr", "", "server address")
	startCmd.PersistentFlags().StringVar(&StatAddr, "statAddr", "", "stat address")
	return startCmd
}

func AddStartCommand() {
	rootCmd.AddCommand(startCmd)
}
