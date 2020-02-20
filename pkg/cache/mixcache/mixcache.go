package mixcache

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gomodule/redigo/redis"
	"github.com/i2eco/muses/pkg/cache/mixcache/gocache"
	"github.com/i2eco/muses/pkg/cache/mixcache/leveldb"
	mredis "github.com/i2eco/muses/pkg/cache/mixcache/redis"
	"github.com/i2eco/muses/pkg/cache/mixcache/standard"
	"github.com/i2eco/muses/pkg/common"
	"sync"
)

type Error string

func (err Error) Error() string { return string(err) }

var ErrNil = errors.New("mixcache: nil returned")

var defaultCaller = &callerStore{
	Name: common.ModMixCacheName,
}

type callerStore struct {
	Name         string
	IsBackground bool
	caller       sync.Map
	cfg          Cfg
}

type Client struct {
	standard.MixCache
	cfg CallerCfg
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
	for name, cfg := range c.cfg.Muses.Mixcache {
		db, err := provider(cfg)
		if err != nil {
			return err
		}
		c := &Client{
			db,
			cfg,
		}
		defaultCaller.caller.Store(name, c)
	}
	return nil
}

func provider(cfg CallerCfg) (client standard.MixCache, err error) {
	if cfg.Mode == "redis" {
		dialOptions := []redis.DialOption{
			redis.DialConnectTimeout(cfg.ConnectTimeout.Duration),
			redis.DialReadTimeout(cfg.ReadTimeout.Duration),
			redis.DialWriteTimeout(cfg.WriteTimeout.Duration),
			redis.DialDatabase(cfg.DB),
			redis.DialPassword(cfg.Password),
		}
		client, err = mredis.NewMixCache(cfg.Addr, cfg.Debug, cfg.MaxActive, cfg.IdleTimeout.Duration, cfg.Wait, dialOptions...)
	} else if cfg.Mode == "memory" {
		client, err = gocache.NewMixCache()
	} else if cfg.Mode == "leveldb" {
		client, err = leveldb.NewMixCache(cfg.LevelDbPath)
	} else {
		err = errors.New("oss mode not exist")
		return
	}
	return
}

func String(reply interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}
	switch reply := reply.(type) {
	case []byte:
		return string(reply), nil
	case string:
		return reply, nil
	case nil:
		return "", ErrNil
	case Error:
		return "", reply
	}
	return "", fmt.Errorf("redigo: unexpected type for String, got type %T", reply)
}
