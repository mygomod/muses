package ginsession

type Cfg struct {
	Muses struct {
		Session struct {
			CallerCfg `toml:"gin"`
		} `toml:"session"`
	} `toml:"muses"`
}

type CallerCfg struct {
	Mode string // redis,memstore

	Name     string
	Size     int
	Network  string
	Addr     string
	Pwd      string
	Keypairs string
}
