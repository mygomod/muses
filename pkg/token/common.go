package token

import (
	"github.com/goecology/muses/pkg/cache/redis"
	"github.com/goecology/muses/pkg/database/mysql"
	"github.com/goecology/muses/pkg/logger"
)

type Cfg struct {
	Muses struct {
		Token map[string]CallerCfg `toml:"token"`
	} `toml:"muses"`
}

type MysqlCallerCfg mysql.CallerCfg
type RedisCallerCfg redis.CallerCfg
type LoggerCallerCfg logger.CallerCfg

type CallerCfg struct {
	Mode string

	LoggerRef string
	LoggerCallerCfg

	MysqlRef string
	MysqlCallerCfg

	RedisRef string
	RedisCallerCfg
}
