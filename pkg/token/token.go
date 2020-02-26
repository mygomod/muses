package token

import (
	"errors"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"


	mysqlToken "github.com/i2eco/muses/pkg/token/mysql"
	redis2 "github.com/i2eco/muses/pkg/token/redis"
	"github.com/i2eco/muses/pkg/token/standard"

	"github.com/i2eco/muses/pkg/cache/redis"
	"github.com/i2eco/muses/pkg/common"
	"github.com/i2eco/muses/pkg/database/mysql"
	"github.com/i2eco/muses/pkg/logger"
)

var defaultCallerStore = &callerStore{
	Name: common.ModTokenName,
}

type callerStore struct {
	Name   string
	caller sync.Map
	cfg    Cfg
}

type Client struct {
	standard.TokenAccessor
	cfg CallerCfg
}

func Register() common.Caller {
	return defaultCallerStore
}

func Caller(name string) *Client {
	obj, ok := defaultCallerStore.caller.Load(name)
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
	for name, cfg := range c.cfg.Muses.Token {
		accessor, err := provider(cfg)
		if err != nil {
			return err
		}
		c := &Client{
			accessor,
			cfg,
		}
		defaultCallerStore.caller.Store(name, c)
	}
	return nil
}

func provider(cfg CallerCfg) (client standard.TokenAccessor, err error) {
	var loggerClient *logger.Client

	// 如果没有引用的logger，就创建一个
	if len(cfg.LoggerRef) > 0 {
		loggerClient = logger.Caller(cfg.LoggerRef)
	} else {
		loggerClient = logger.Provider(logger.CallerCfg(cfg.Logger))
	}

	if cfg.Mode == "mysql" {
		return createMysqlAccessor(cfg, loggerClient)
	} else if cfg.Mode == "redis" {
		return createRedisAccessor(cfg, loggerClient)
	} else {
		return nil, errors.New("The token's mode must be redis or mysql: " + cfg.Mode)
	}
}

func createMysqlAccessor(cfg CallerCfg, loggerClient *logger.Client) (accessor standard.TokenAccessor, err error) {
	var db *gorm.DB
	if len(cfg.MysqlRef) > 0 {
		db = mysql.Caller(cfg.MysqlRef)
	} else {
		db, err = mysql.Provider(mysql.CallerCfg(cfg.Mysql))
		if err != nil {
			return
		}
	}
	return mysqlToken.InitTokenAccessor(loggerClient, db), nil
}

func createRedisAccessor(cfg CallerCfg, loggerClient *logger.Client) (standard.TokenAccessor, error) {
	var redisClient *redis.Client
	if len(cfg.RedisRef) > 0 {
		redisClient = redis.Caller(cfg.RedisRef)
	} else {
		redisClient = redis.Provider(redis.CallerCfg(cfg.Redis))
	}

	return redis2.InitRedisTokenAccessor(loggerClient, redisClient), nil
}
