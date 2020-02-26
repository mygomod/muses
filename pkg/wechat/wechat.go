package wechat

import (
	"github.com/BurntSushi/toml"
	"github.com/go-resty/resty/v2"
	"github.com/i2eco/muses/pkg/common"
	"github.com/i2eco/muses/pkg/wechat/cache"
	"github.com/i2eco/muses/pkg/wechat/context"
	"github.com/i2eco/muses/pkg/wechat/miniprogram"
	"sync"
)

var defaultCaller = &callerStore{
	Name: common.ModWechatName,
}

type callerStore struct {
	Name   string
	cfg    Cfg
	caller *Client
}

type Client struct {
	RestyClient *resty.Client
	Cache       cache.Cache
	Context     *context.Context
}

func Register() common.Caller {
	return defaultCaller
}

func Caller() *Client {
	return defaultCaller.caller
}

func (c *callerStore) InitCfg(cfg []byte) error {
	if err := toml.Unmarshal(cfg, &c.cfg); err != nil {
		return err
	}
	return nil
}

func (c *callerStore) InitCaller() error {
	db, err := provider(c.cfg.Muses.Wechat)
	if err != nil {
		return err
	}
	defaultCaller.caller = db
	return nil
}

func provider(cfg CallerCfg) (resp *Client, err error) {
	resp = &Client{}
	if cfg.Mode == "redis" {
		resp.Cache = cache.NewRedis(&cache.RedisOpts{
			Debug:     cfg.Debug,
			Addr:      cfg.Addr,
			Password:  cfg.Password,
			Database:  cfg.Database,
			MaxIdle:   cfg.MaxIdle,
			MaxActive: cfg.MaxActive,
		})
	} else {
		panic("mode not support")
	}
	ctx := new(context.Context)
	ctx.AppID = cfg.AppID
	ctx.AppSecret = cfg.AppSecret
	ctx.Token = cfg.Token
	ctx.EncodingAESKey = cfg.EncodingAESKey
	ctx.PayMchID = cfg.PayMchID
	ctx.PayKey = cfg.PayKey
	ctx.PayNotifyURL = cfg.PayNotifyURL
	ctx.Cache = resp.Cache
	ctx.SetAccessTokenLock(new(sync.RWMutex))
	ctx.SetJsAPITicketLock(new(sync.RWMutex))
	ctx.RestyClient = resty.New().SetDebug(cfg.Debug)
	resp.Context = ctx
	return
}

// GetMiniProgram 获取小程序的实例
func (wc *Client) GetMiniProgram() *miniprogram.MiniProgram {
	return miniprogram.NewMiniProgram(wc.Context)
}
