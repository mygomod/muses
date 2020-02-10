package oss

type Cfg struct {
	Muses struct {
		Oss map[string]CallerCfg `toml:"oss"`
	} `toml:"muses"`
}

type CallerCfg struct {
	Debug           bool
	Mode            string
	Addr            string
	AccessKeyID     string
	AccessKeySecret string
	CdnName         string
	OssBucket       string
	FileBucket      string
	IsDeleteSrcPath bool
}
