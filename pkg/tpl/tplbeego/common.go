package tplbeego

type Cfg struct {
	Muses CallerMuses `toml:"muses"`
}

type CallerMuses struct {
	Tpl CallerTpl `toml:"tpl"`
}

type CallerTpl struct {
	TplBeego CallerCfg `toml:"tplBeego"`
}

type CallerCfg struct {
	Debug         bool
	TplExt        string
	ViewPath      string
	TemplateLeft  string
	TemplateRight string
}
