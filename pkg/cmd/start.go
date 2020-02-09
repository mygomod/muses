package cmd

import (
	"github.com/goecology/muses/pkg/system"
	"github.com/spf13/cobra"
)

var ConfigPath string
var Addr string

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
	startCmd.PersistentFlags().StringVarP(&Addr, "addr", "", "", "address")
	return startCmd
}

func AddStartCommand() {
	rootCmd.AddCommand(startCmd)
}
