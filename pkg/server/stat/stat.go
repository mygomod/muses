package stat

import (
	"fmt"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"

	"github.com/i2eco/muses/pkg/cmd"
	"github.com/i2eco/muses/pkg/common"
	"github.com/i2eco/muses/pkg/logger"
	"github.com/i2eco/muses/pkg/registry"

	"github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
)

var defaultCaller = &callerStore{
	Name: common.ModStatName,
}

type callerStore struct {
	Name   string
	caller *Client
	cfg    Cfg
}

type Client struct {
	*http.Server
}

func Register() common.Caller {
	return defaultCaller
}

func (c *callerStore) InitCfg(cfg []byte) error {
	if err := toml.Unmarshal(cfg, &c.cfg); err != nil {
		return err
	}
	return nil
}

func (c *callerStore) InitCaller() error {
	addr := c.cfg.Muses.Server.Stat.Addr
	// 如果存在命令行的addr，覆盖配置里的addr
	if cmd.StatAddr != "" {
		addr = cmd.StatAddr
	}

	serverStats := &http.Server{
		Addr:         addr,
		Handler:      initStat(),
		ReadTimeout:  c.cfg.Muses.Server.Stat.ReadTimeout.Duration,
		WriteTimeout: c.cfg.Muses.Server.Stat.WriteTimeout.Duration,
	}

	go func() {
		defer func() {
			serverStats.Close()
		}()
		if err := serverStats.ListenAndServe(); err != nil {
			fmt.Println("err", c.cfg.Muses.Server)
			logger.DefaultLogger().Error("ServerApi err", zap.String("err", err.Error()))
		}
	}()
	c.caller = &Client{serverStats}
	// registry the ServerStatAddr
	err := registry.RegisterServerStatAddr(addr)
	if err != nil {
		logger.DefaultLogger().Error("Registry server state address failed, err",
			zap.String("err", err.Error()))
	}
	return nil
}

func initStat() http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())
	p := ginprometheus.NewPrometheus(common.MetricPrefix)
	p.Use(r)
	r.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome ServerApi Stats",
			},
		)
	})

	return r
}
