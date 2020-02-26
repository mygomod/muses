package wechat

type Cfg struct {
	Muses struct {
		Wechat CallerCfg `toml:"wechat"`
	} `toml:"muses"`
}

type CallerCfg struct {
	Mode        string `toml:"mode"`
	Debug       bool   `toml:"debug"`
	Addr        string `toml:"addr" json:"addr"`
	Password    string `toml:"password" json:"password"`
	Database    int    `toml:"db" json:"database"`
	MaxIdle     int    `toml:"maxIdle" json:"max_idle"`
	MaxActive   int    `toml:"maxActive" json:"max_active"`
	IdleTimeout int32  `toml:"idleTimeout" json:"idle_timeout"` //second

	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string
	PayMchID       string
	PayKey         string
	PayNotifyURL   string
}
