package oss

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/mygomod/muses/pkg/common"
	"github.com/mygomod/muses/pkg/oss/alioss"
	"github.com/mygomod/muses/pkg/oss/file"
	"github.com/mygomod/muses/pkg/oss/standard"
	"github.com/satori/uuid"
	"strings"
	"sync"
	"time"
)

var defaultCaller = &callerStore{
	Name: common.ModOssName,
}

type callerStore struct {
	Name         string
	IsBackground bool
	caller       sync.Map
	cfg          Cfg
}

type Client struct {
	standard.Oss
	cfg CallerCfg
}

func Register() common.Caller {
	return defaultCaller
}

func Caller(name string) *Client {
	obj, ok := defaultCaller.caller.Load(name)
	if !ok {
		return nil
	}
	return obj.(*Client)
}

func (c *callerStore) InitCfg(cfg []byte) error {
	if err := toml.Unmarshal(cfg, &c.cfg); err != nil {
		return err
	}
	return nil
}

func (c *callerStore) InitCaller() error {
	for name, cfg := range c.cfg.Muses.Oss {
		db, err := provider(cfg)
		if err != nil {
			return err
		}
		c := &Client{
			db,
			cfg,
		}
		defaultCaller.caller.Store(name, c)
	}
	return nil
}

func provider(cfg CallerCfg) (client standard.Oss, err error) {
	if cfg.Mode == "alioss" {
		client, err = alioss.NewOss(cfg.Addr, cfg.AccessKeyID, cfg.AccessKeySecret, cfg.OssBucket, cfg.IsDeleteSrcPath)
	} else if cfg.Mode == "file" {
		client, err = file.NewOss(cfg.CdnName, cfg.FileBucket, cfg.IsDeleteSrcPath)
	} else {
		err = errors.New("oss mode not exist")
		return
	}
	return
}

func (c *Client) ShowImg(img string, style ...string) (url string) {
	if strings.HasPrefix(img, "https://") || strings.HasPrefix(img, "http://") {
		return img
	}
	img = strings.TrimLeft(img, "./")
	switch c.cfg.Mode {
	case "alioss":
		s := ""
		if len(style) > 0 && strings.TrimSpace(style[0]) != "" {
			s = "/" + style[0]
		}
		url = img + s
	case "file":
		url = img
	}
	url = c.cfg.CdnName + url
	return
}

func (c *Client) ShowImgArr(imgs []string, style ...string) (urlArr []string) {
	urlArr = make([]string, 0)
	for _, img := range imgs {
		urlArr = append(urlArr, c.ShowImg(img, style...))
	}
	return
}

func (c *Client) GenerateKey(prefix string) string {
	month := time.Now().Format("200601")
	// Reossv2上传报错：Thespecifiedobjectisnotvalid.
	// object路径开头不能与“/”
	return prefix + "/" + month + "/" + strings.ReplaceAll(uuid.NewV4().String(), "-", "") + ".jpg"
}
