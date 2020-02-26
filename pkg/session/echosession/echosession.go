package echosession

import (
	"github.com/BurntSushi/toml"
	"github.com/ipfans/echo-session"
	"github.com/labstack/echo"
	session "github.com/sevenNt/echo-session"

	"github.com/i2eco/muses/pkg/common"
)

var defaultCaller = &callerStore{
	Name: common.ModEchoSessionName,
}

type callerStore struct {
	Name         string
	IsBackground bool
	caller       echo.MiddlewareFunc
	cfg          Cfg
}

func Register() common.Caller {
	return defaultCaller
}

func Caller() echo.MiddlewareFunc {
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

func provider(cfg CallerCfg) (s echo.MiddlewareFunc, err error) {
	store, err := session.NewRedisStore(cfg.Size, cfg.Network, cfg.Addr, cfg.Pwd, []byte(cfg.Keypairs))
	s = session.Sessions(cfg.Name, store)
	return
}
