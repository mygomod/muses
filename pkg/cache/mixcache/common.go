package mixcache

import "github.com/i2eco/muses/pkg/common"

type Cfg struct {
	Muses struct {
		Mixcache map[string]CallerCfg `toml:"mixcache"`
	} `toml:"muses"`
}

type CallerCfg struct {
	Debug          bool
	Mode           string
	LevelDbPath    string
	Network        string // tcp
	Addr           string // 127.0.0.1:6379
	DB             int
	Password       string
	ConnectTimeout common.Duration
	ReadTimeout    common.Duration
	WriteTimeout   common.Duration
	MaxIdle        int
	MaxActive      int
	IdleTimeout    common.Duration
	Wait           bool
}
