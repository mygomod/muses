package redis

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gomodule/redigo/redis"

	"github.com/goecology/muses/pkg/common"
)

var defaultCaller = &callerStore{
	Name: common.ModRedisName,
}

type callerStore struct {
	Name         string
	IsBackground bool
	caller       sync.Map
	cfg          Cfg
}

type Client struct {
	pool *redis.Pool
}

func Register() common.Caller {
	return defaultCaller
}

func Caller(name string) *Client {
	obj, ok := defaultCaller.caller.Load(name)
	if !ok {
		return nil
	}
	return obj.(*Client)
}

func (c *callerStore) InitCfg(cfg []byte) error {
	if err := toml.Unmarshal(cfg, &c.cfg); err != nil {
		return err
	}
	return nil
}

func (c *callerStore) InitCaller() error {
	for name, cfg := range c.cfg.Muses.Redis {
		db := Provider(cfg)
		defaultCaller.caller.Store(name, db)
	}
	return nil
}

func Provider(cfg CallerCfg) (resp *Client) {
	dialOptions := []redis.DialOption{
		redis.DialConnectTimeout(cfg.ConnectTimeout.Duration),
		redis.DialReadTimeout(cfg.ReadTimeout.Duration),
		redis.DialWriteTimeout(cfg.WriteTimeout.Duration),
		redis.DialDatabase(cfg.DB),
		redis.DialPassword(cfg.Password),
	}

	resp = &Client{
		&redis.Pool{
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", cfg.Addr, dialOptions...)
				if err != nil {
					return nil, err
				}

				if cfg.Debug {
					return redis.NewLoggingConn(c, log.New(os.Stderr, "", log.LstdFlags), "redis"), nil
				}
				return c, nil
			},
			// Use the TestOnBorrow function to check the health of an idle connection
			// before the connection is returned to the application. This example PINGs
			// connections that have been idle more than a minute:
			//
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
			MaxIdle:     cfg.MaxActive,
			MaxActive:   cfg.MaxActive,
			IdleTimeout: cfg.IdleTimeout.Duration,
			Wait:        cfg.Wait, // wait until getting connection from pool
		},
	}
	return
}
