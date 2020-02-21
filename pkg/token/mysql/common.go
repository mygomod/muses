package mysql

type Cfg struct {
	Muses struct {
		Token struct {
			JWT struct {
				CallerCfg `toml: "mysql"`
			} `toml: "jwt"`
		} `toml:"token"`
	} `toml: "muses"`
}

type CallerCfg struct {
	logger string // logger client
	client string // redis client name
}
