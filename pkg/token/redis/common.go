package redis

type Cfg struct {
	Muses struct{
		Token struct{
			JWT struct{
				CallerCfg `toml:"redis"`
			} `toml:"jwt"`
		} `toml:"token"`
	} `toml:"muses"`
}

type CallerCfg struct {
	logger string // logger client
	client string // redis client name
}
