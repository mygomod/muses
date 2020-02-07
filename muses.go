package muses

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fvbock/endless"
	ogin "github.com/gin-gonic/gin"
	"github.com/goecology/muses/pkg/app"
	"github.com/goecology/muses/pkg/cmd"
	"github.com/goecology/muses/pkg/common"
	"github.com/goecology/muses/pkg/logger"
	"github.com/goecology/muses/pkg/prom"
	"github.com/goecology/muses/pkg/server/gin"
	"github.com/goecology/muses/pkg/system"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"syscall"
)

type Muses struct {
	cfgByte     []byte
	callers     []common.Caller
	isSetConfig bool
	filePath    string
	preRun      []common.PreRunFunc
	ext         string
	err         error
	router      func() *ogin.Engine
}

func Container(callerFuncs ...common.CallerFunc) (muses *Muses) {
	allCallers := []common.CallerFunc{app.Register, logger.Register, prom.Register}
	allCallers = append(allCallers, callerFuncs...)
	muses = &Muses{}
	callers, err := sortCallers(allCallers)
	if err != nil {
		muses.err = err
		return
	}
	muses.callers = callers
	// 初始化 启动信息
	system.InitRunInfo()
	muses.showVersion()
	return
}

func (m *Muses) SetRouter(router func() *ogin.Engine) *Muses {
	if m.err != nil {
		return m
	}
	var err error
	m.router = router
	cmd.SetStartFn(m.startFn)
	for _, caller := range m.callers {
		name := getCallerName(caller)
		// 说明启动cmd配置，那么就不需要在setConfig
		if name == common.ModCmdName {
			// 启动配置选项
			err = caller.InitCfg([]byte{})
			if err != nil {
				m.err = err
				return m
			}
			m.filePath = common.CmdConfigPath
			err = isPathExist(m.filePath)
			if err != nil {
				m.err = err
				return m
			}

			ext := filepath.Ext(m.filePath)

			if len(ext) <= 1 {
				m.err = errors.New("config file ext is error")
				return m
			}
			m.ext = ext[1:]
			m.SetCfg(m.filePath)
		}
	}
	return m
}

func (m *Muses) SetCfg(cfg interface{}) *Muses {
	if m.err != nil {
		return m
	}
	var err error
	var cfgByte []byte
	switch cfg.(type) {
	case string:
		cfgByte, err = parseFile(cfg.(string))
		if err != nil {
			m.err = err
			return m
		}
		m.filePath = cfg.(string)
	case []byte:
		cfgByte = cfg.([]byte)
	default:
		m.err = fmt.Errorf("type is error %s", cfg)
		return m
	}
	m.cfgByte = cfgByte
	m.isSetConfig = true
	m.initViper()
	return m
}

func (m *Muses) PreRun(f ...common.PreRunFunc) *Muses {
	m.preRun = f
	return m
}

func (m *Muses) Run() (err error) {
	fmt.Println(system.BuildInfo.LongForm())
	if !m.isSetConfig {
		err = errors.New("config file is not setting")
		return
	}

	for _, caller := range m.callers {
		name := getCallerName(caller)
		// 说明是最后一个启动指令
		if name == common.ModCmdName {
			// 运行前置指令
			for _, f := range m.preRun {
				err = f()
				if err != nil {
					return
				}
			}
			m.printInfo("module", name, "init caller start")
			if err = caller.InitCaller(); err != nil {
				m.printInfo("module", name, "init caller error")
				return
			}
			m.printInfo("module", name, "init caller ok")
		} else {
			m.printInfo("module", name, "cfg start")
			if err = caller.InitCfg(m.cfgByte); err != nil {
				m.printInfo("module", name, "init config error")
				return
			}
			m.printInfo("module", name, "cfg end")
			m.printInfo("module", name, "init caller start")
			if err = caller.InitCaller(); err != nil {
				m.printInfo("module", name, "init caller error")
				return
			}
			m.printInfo("module", name, "init caller ok")
		}

	}
	return
}

func (m *Muses) initViper() {
	rBytes := bytes.NewReader(m.cfgByte)
	viper.SetConfigType(m.ext)
	viper.AutomaticEnv() // read in environment variables that match
	//viper.Debug()
	err := viper.ReadConfig(rBytes)
	if err != nil {
		m.printInfo("Using config file:", viper.ConfigFileUsed())
	}

}

func isPathExist(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	return err
}

func (m *Muses) startFn(cmd *cobra.Command, args []string) {
	// 主服务器
	endless.DefaultReadTimeOut = gin.Config().Muses.Server.Gin.ReadTimeout.Duration
	endless.DefaultWriteTimeOut = gin.Config().Muses.Server.Gin.WriteTimeout.Duration
	endless.DefaultMaxHeaderBytes = 100000000000000
	server := endless.NewServer(gin.Config().Muses.Server.Gin.Addr, m.router())
	server.BeforeBegin = func(add string) {
		logger.DefaultLogger().Info(fmt.Sprintf("Actual pid is %d", syscall.Getpid()))
	}

	if err := server.ListenAndServe(); err != nil {
		logger.DefaultLogger().Error("Server err", zap.String("err", err.Error()))
	}
}

func isVersion() bool {
	if len(os.Args) > 1 && "version" == os.Args[1] {
		return true
	}
	return false
}

func (m *Muses) showVersion() {
	if isVersion() {
		for _, caller := range m.callers {
			name := getCallerName(caller)
			// 说明是最后一个启动指令
			if name == common.ModCmdName {
				_ = caller.InitCaller()
				os.Exit(0)
			}
		}
	}
}

func (m *Muses) printInfo(info ...interface{}) {
	fmt.Println(info)
}
