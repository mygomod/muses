package cmd

type Cfg struct {
	Muses struct {
		Cmd CallerCfg `toml:"app"`
	} `toml:"muses"`
}

type CallerCfg struct{}
