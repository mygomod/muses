package ginsession

import (
	"github.com/BurntSushi/toml"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/goecology/muses/pkg/common"
)

var defaultCaller = &callerStore{
	Name: common.ModGinSessionName,
}

type callerStore struct {
	Name         string
	IsNecessary  bool
	IsBackground bool
	caller       gin.HandlerFunc
	cfg          Cfg
}

func Register() common.Caller {
	return defaultCaller
}

func Caller() gin.HandlerFunc {
	return defaultCaller.caller
}

func (c *callerStore) InitCfg(cfg []byte) error {
	if err := toml.Unmarshal(cfg, &c.cfg); err != nil {
		return err
	}
	return nil
}

func (c *callerStore) InitCaller() error {
	db, err := provider(c.cfg.Muses.Session.CallerCfg)
	if err != nil {
		return err
	}
	c.caller = db
	return nil
}

func provider(cfg CallerCfg) (session gin.HandlerFunc, err error) {
	var store redis.Store
	store, err = redis.NewStore(cfg.Size, cfg.Network, cfg.Addr, cfg.Pwd, []byte(cfg.Keypairs))
	session = sessions.Sessions(cfg.Name, store)
	return
}
