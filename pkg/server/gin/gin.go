package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/goecology/muses/pkg/common"
	"github.com/goecology/muses/pkg/server/gin/plugins/ginzap"
	"time"

	"github.com/BurntSushi/toml"
)

var defaultCaller = &callerStore{
	Name: common.ModGinName,
}

type callerStore struct {
	Name   string
	caller *gin.Engine
	cfg    Cfg
}

func Register() common.Caller {
	return defaultCaller
}

func Caller() *gin.Engine {
	return defaultCaller.caller
}

func Config() Cfg {
	return defaultCaller.cfg
}

func (c *callerStore) InitCfg(cfg []byte) error {
	if err := toml.Unmarshal(cfg, &c.cfg); err != nil {
		return err
	}
	return nil
}

func (c *callerStore) InitCaller() error {
	gin.SetMode(c.cfg.Muses.Server.Gin.Mode)

	r := gin.New()
	if c.cfg.Muses.Server.Gin.EnabledLogger {
		r.Use(ginzap.Ginzap(time.RFC3339, true, c.cfg.Muses.Server.Gin.EnabledMetric))
	}

	if c.cfg.Muses.Server.Gin.EnabledRecovery {
		r.Use(ginzap.RecoveryWithZap(true))
	}

	c.caller = r
	return nil
}
