package tplbeego

import (
	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/goecology/muses/pkg/common"
)

var defaultCaller = &callerStore{
	Name: common.ModTplBeegoName,
	cfg: Cfg{
		Muses: CallerMuses{
			Tpl: CallerTpl{
				TplBeego: CallerCfg{
					Debug:         false,
					TplExt:        "tpl",
					ViewPath:      "views",
					TemplateLeft:  "{{",
					TemplateRight: "}}",
				},
			},
		},
	},
}

type callerStore struct {
	Name string
	cfg  Cfg
}

func Register() common.Caller {
	return defaultCaller
}

func Caller() (obj *Tmpl, err error) {
	obj, err = provider(defaultCaller.cfg.Muses.Tpl.TplBeego)
	if err != nil {
		return
	}
	return
}

func Config() CallerCfg {
	return defaultCaller.cfg.Muses.Tpl.TplBeego
}

func (c *callerStore) InitCfg(cfg []byte) error {
	if err := toml.Unmarshal(cfg, &c.cfg); err != nil {
		return err
	}
	return nil
}

func (c *callerStore) InitCaller() error {
	if err := AddViewPath(c.cfg.Muses.Tpl.TplBeego.ViewPath); err != nil {
		return err
	}
	return nil
}

func provider(cfg CallerCfg) (resp *Tmpl, err error) {
	obj := &Tmpl{}
	obj.Init(cfg.TplExt, cfg.ViewPath)

	return obj, err
}
