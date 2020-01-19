package stat

import (
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/goecology/muses/pkg/common"
	"github.com/goecology/muses/pkg/logger"
	"go.uber.org/zap"
	"net/http"
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
	serverStats := &http.Server{
		Addr:         c.cfg.Muses.Server.Stat.Addr,
		Handler:      initStat(),
		ReadTimeout:  c.cfg.Muses.Server.Stat.ReadTimeout.Duration,
		WriteTimeout: c.cfg.Muses.Server.Stat.WriteTimeout.Duration,
	}

	go func() {
		defer func() {
			serverStats.Close()
		}()
		if err := serverStats.ListenAndServe(); err != nil {
			logger.DefaultLogger().Error("ServerApi err", zap.String("err", err.Error()))
		}
	}()
	c.caller = &Client{serverStats}

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
