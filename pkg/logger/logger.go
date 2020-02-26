package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/i2eco/muses/pkg/common"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultCaller = &callerStore{
	Name: common.ModLoggerName,
}

type callerStore struct {
	Name         string
	IsBackground bool
	caller       sync.Map
	cfg          Cfg
}

type Client struct {
	*zap.Logger
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

// 初始化做了判断，肯定存在默认配置
func DefaultLogger() *Client {
	var logClient *Client
	// 如果设置了系统日志，就返回系统日志
	obj, ok := defaultCaller.caller.Load(common.SystemLogger)
	if !ok {
		// 如果没有系统日志，那么就返回用户设置的第一个日志
		defaultCaller.caller.Range(func(key, value interface{}) bool {
			logClient = value.(*Client)
			return false
		})
	} else {
		logClient = obj.(*Client)
	}

	// 如果log client 不存在，提示用户配置里需要设置日志配置
	if logClient == nil {
		panic("please set logger config")
	}
	return logClient
}

func (c *callerStore) InitCfg(cfg []byte) error {
	if err := toml.Unmarshal(cfg, &c.cfg); err != nil {
		return err
	}
	return nil
}

func (c *callerStore) InitCaller() error {
	for name, cfg := range c.cfg.Muses.Logger {
		db := Provider(cfg)
		defaultCaller.caller.Store(name, db)
	}
	return nil
}

func Provider(cfg CallerCfg) (db *Client) {
	var js string
	if cfg.Debug {
		js = fmt.Sprintf(`{
      "level": "%s",
      "encoding": "json",
      "outputPaths": ["stdout"],
      "errorOutputPaths": ["stdout"]
      }`, cfg.Level)
	} else {
		js = fmt.Sprintf(`{
      "level": "%s",
      "encoding": "json",
      "outputPaths": ["%s"],
      "errorOutputPaths": ["%s"]
      }`, cfg.Level, cfg.Path, cfg.Path)
	}

	var zcfg zap.Config
	if err := json.Unmarshal([]byte(js), &zcfg); err != nil {
		panic(err)
	}
	zcfg.EncoderConfig = zap.NewProductionEncoderConfig()
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var err error
	l, err := zcfg.Build()
	if err != nil {
		log.Fatal("init logger error: ", err)
	}
	return &Client{l}
}
