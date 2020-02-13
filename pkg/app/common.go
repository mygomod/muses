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

const separator = "-"

// AppKey will return a string representing current APP.
// The result looks like Name-Version-Env
func (cfg CallerCfg) AppKey() string {
	result := cfg.Name
	if len(cfg.Version) > 0 {
		result = result + separator + cfg.Version
	}

	if len(cfg.Env) > 0 {
		result = result + separator + cfg.Env
	}
	return result
}
