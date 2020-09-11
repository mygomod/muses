package token

import (
	"github.com/mygomod/muses/pkg/cache/redis"
	"github.com/mygomod/muses/pkg/database/mysql"
	"github.com/mygomod/muses/pkg/logger"
)

type Cfg struct {
	Muses struct {
		Token map[string]CallerCfg `toml:"token"`
	} `toml:"muses"`
}

type MysqlCallerCfg mysql.CallerCfg
type RedisCallerCfg redis.CallerCfg
type LoggerCallerCfg logger.CallerCfg

// CallerCfg是token的配置。
// 需要注意的是，XXXRef是指指定了一个已经在配置文件里面的Caller。
// 比如说，你已经设置了一个mysql数据库myDB用于存储数据，而你又希望同时使用该数据库来存放token的数据，
// 那么你只需将MysqlRef设置为myDB。
// 如果你没有指定Ref，那么token会在初始化的时候依据配置来重新创建一个
type CallerCfg struct {
	Mode string

	LoggerRef string
	Logger    LoggerCallerCfg `toml:"logger"`

	MysqlRef string
	Mysql    MysqlCallerCfg `toml:"mysql"`

	RedisRef string
	Redis    RedisCallerCfg `toml:"redis"`
}
