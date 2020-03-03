package github

type Cfg struct {
	Muses struct {
		OpenGithub CallerCfg `toml:"openGithub"`
	} `toml:"muses"`
}

// CallerCfg是Github的配置。
type CallerCfg struct {
	Debug        bool
	Mode         string
	ClientID     string
	ClientSecret string
	RedirectURL  string
}
