package muses

import (
	"fmt"
	"github.com/goecology/muses/pkg/app"
	"github.com/goecology/muses/pkg/common"
	"github.com/goecology/muses/pkg/logger"
	"github.com/goecology/muses/pkg/prom"
	"github.com/goecology/muses/pkg/system"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var filePath string

func Container(cfg interface{}, callerFuncs ...common.CallerFunc) (err error) {
	var cfgByte []byte
	switch cfg.(type) {
	case string:
		cfgByte, err = parseFile(cfg.(string))
		if err != nil {
			return
		}
		filePath = cfg.(string)
	case []byte:
		cfgByte = cfg.([]byte)
	default:
		return fmt.Errorf("type is error %s", cfg)
	}

	// 初始化 viper 配置信息
	initViper()
	// 初始化 启动信息
	system.InitRunInfo()
	fmt.Println(system.BuildInfo.LongForm())
	fmt.Println(system.RunInfo.LongForm())

	allCallers := []common.CallerFunc{app.Register, logger.Register, prom.Register}
	allCallers = append(allCallers, callerFuncs...)

	callers, err := sortCallers(allCallers)
	if err != nil {
		return
	}

	for _, caller := range callers {
		name := getCallerName(caller)
		fmt.Println("module", name, "start")
		if err = caller.InitCfg(cfgByte); err != nil {
			fmt.Println("module", name, "init config error")
			return
		}
		fmt.Println("module", name, "init config ok")
		if err = caller.InitCaller(); err != nil {
			fmt.Println("module", name, "init caller error")
			return
		}
		fmt.Println("module", name, "init caller ok")
		fmt.Println("module", name, "end")
	}
	return nil
}

// todo
func initViper() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigFile(filePath)
	viper.AutomaticEnv() // read in environment variables that match
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
