package mongo

type Cfg struct {
	Muses struct {
		Mongo map[string]CallerCfg `toml:"mongo"`
	} `toml:"muses"`
}

type CallerCfg struct {
	Debug bool

	URL      string
	Database string
}
