package echosession

type Cfg struct {
	Muses struct {
		Session struct {
			CallerCfg `toml:"echo"`
		} `toml:"session"`
	} `toml:"muses"`
}

type CallerCfg struct {
	Name     string
	Size     int
	Network  string
	Addr     string
	Pwd      string
	Keypairs string
}
