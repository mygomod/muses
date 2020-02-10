package etcd

import (
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/coreos/etcd/clientv3"
	"github.com/prometheus/common/log"

	"github.com/goecology/muses/pkg/common"
	"github.com/goecology/muses/pkg/registry"
)

var defaultCaller = &callerStore{
	Name:   common.ModEtcd,
}

type callerStore struct {
	Name string
	caller sync.Map
	cfg Cfg
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
		db, err := clientv3.New(cfg)
		if err != nil {
			return err
		}
		defaultCaller.caller.Store(name, db)
	}
	return nil
}


type EtcdRegistry struct {
	client *clientv3.Client
}

func (e EtcdRegistry) Register(key string, value interface{}) {

}

var (
	etcdRegistryInitOnce  sync.Once
	etcdRegistrySingleton *EtcdRegistry
)

// GetEtcdRegistry will return a EtcdRegistry instance. If it's not intialized,
// we will create one.
// This method is thead-safe.
// return nil if we could not initialize the instance
func GetEtcdRegistry() registry.Registry {
	etcdRegistryInitOnce.Do(func() {
		client, err := clientv3.New(clientv3.Config{

		})
		if err != nil {
			log.Errorf("Could not create the etcd client, the error is: %v", err)
			return
		}
		etcdRegistrySingleton = &EtcdRegistry{
			client: client,
		}
	})
	return etcdRegistrySingleton
}
