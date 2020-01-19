package gin

import "github.com/goecology/muses/pkg/common"

type Cfg struct {
	Muses struct {
		Server struct {
			Gin CallerCfg `toml:"gin"`
		} `toml:"server"`
	} `toml:"muses"`
}

type CallerCfg struct {
	Mode            string
	Addr            string
	ReadTimeout     common.Duration
	WriteTimeout    common.Duration
	MaxHeaderBytes  int
	EnabledRecovery bool
	EnabledLogger   bool
	EnabledMetric   bool
}
