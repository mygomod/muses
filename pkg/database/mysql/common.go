package mysql

import (
	"github.com/goecology/muses/pkg/common"
)

type Cfg struct {
	Muses struct {
		Mysql map[string]CallerCfg `toml:"mysql"`
	} `toml:"muses"`
}

type CallerCfg struct {
	Debug bool

	Network      string
	Dialect      string
	Addr         string
	Username     string
	Password     string
	Db           string
	Charset      string
	ParseTime    string
	Loc          string
	Timeout      common.Duration
	ReadTimeout  common.Duration
	WriteTimeout common.Duration

	Level           string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime common.Duration
}
