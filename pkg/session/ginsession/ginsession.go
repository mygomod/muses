package ginsession

import (
	"github.com/BurntSushi/toml"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/mygomod/muses/pkg/common"
	"github.com/pkg/errors"
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
	if cfg.Mode == "redis" {
		var store redis.Store
		store, err = redis.NewStore(cfg.Size, cfg.Network, cfg.Addr, cfg.Pwd, []byte(cfg.Keypairs))
		if err != nil {
			return
		}
		session = sessions.Sessions(cfg.Name, store)
	} else if cfg.Mode == "memory" {
		store := memstore.NewStore([]byte(cfg.Keypairs))
		session = sessions.Sessions(cfg.Name, store)
	} else {
		err = errors.New("gin session mode not exist")
		return
	}
	return
}
