package etcd

import (
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/coreos/etcd/clientv3"

	"github.com/goecology/muses/pkg/common"
)

var defaultCaller = &callerStore{
	Name: common.ModEtcdName,
}

type callerStore struct {
	Name   string
	caller sync.Map
	cfg    Cfg
}

func Register() common.Caller {
	return defaultCaller
}

func Caller(name string) *clientv3.Client {
	obj, ok := defaultCaller.caller.Load(name)
	if !ok {
		return nil
	}
	return obj.(*clientv3.Client)
}

func (c *callerStore) InitCfg(cfg []byte) error {
	if err := toml.Unmarshal(cfg, &c.cfg); err != nil {
		return err
	}
	return nil
}

func (c *callerStore) InitCaller() error {
	for name, cfg := range c.cfg.Muses.Etcd {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:            cfg.Endpoints,
			AutoSyncInterval:     cfg.AutoSyncInterval.Duration,
			DialTimeout:          cfg.DialTimeout.Duration,
			DialKeepAliveTime:    cfg.DialKeepAliveTime.Duration,
			DialKeepAliveTimeout: cfg.DialKeepAliveTimeout.Duration,
			MaxCallSendMsgSize:   cfg.MaxCallSendMsgSize,
			MaxCallRecvMsgSize:   cfg.MaxCallRecvMsgSize,
			Username:             cfg.Username,
			Password:             cfg.Password,
			RejectOldCluster:     cfg.RejectOldCluster,
			PermitWithoutStream:  cfg.PermitWithoutStream,
		})
		if err != nil {
			return err
		}
		defaultCaller.caller.Store(name, client)
	}
	return nil
}
