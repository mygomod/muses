package app

type Cfg struct {
	Muses struct {
		App CallerCfg `toml:"app"`
	} `toml:"muses"`
}

type CallerCfg struct {
	Name    string
	Version string
	Env     string
}
