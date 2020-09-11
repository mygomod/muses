package cmd

import (
	"fmt"
	"github.com/mygomod/muses/pkg/system"
	"github.com/spf13/cobra"
)

var short bool

// startCmd represents the hello command
var versionCmd = &cobra.Command{
	Use:  "version",
	Long: `Show version information`,
	Run:  versionFn,
}

func init() {
	versionCmd.PersistentFlags().BoolVarP(&short, "short", "s", short, "Displays a short form of the version information")
	rootCmd.AddCommand(versionCmd)
}

func versionFn(cmd *cobra.Command, args []string) {
	if short {
		fmt.Println(system.BuildInfo)
		return
	} else {
		fmt.Println(system.BuildInfo.LongForm())
		return
	}
}
