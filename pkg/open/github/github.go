package github

import (
	"github.com/BurntSushi/toml"
	"github.com/go-resty/resty/v2"
	"github.com/mygomod/muses/pkg/common"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var defaultCallerStore = &callerStore{
	Name: common.ModOpenGithubName,
}

type callerStore struct {
	Name   string
	cfg    Cfg
	caller *Client
}

type Client struct {
	cfg         CallerCfg
	Oauth       *oauth2.Config
	RestyClient *resty.Client
}

func Register() common.Caller {
	return defaultCallerStore
}

func Caller() *Client {
	return defaultCallerStore.caller
}

func (c *callerStore) InitCfg(cfg []byte) error {
	if err := toml.Unmarshal(cfg, &c.cfg); err != nil {
		return err
	}
	return nil
}

func (c *callerStore) InitCaller() error {
	db, err := provider(c.cfg.Muses.OpenGithub)
	if err != nil {
		return err
	}
	c.caller = db
	return nil
}

func provider(cfg CallerCfg) (resp *Client, err error) {
	resp = &Client{}
	oauth := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  cfg.RedirectURL,
	}
	resp.Oauth = oauth
	resp.RestyClient = resty.New().SetDebug(cfg.Debug)
	return
}
