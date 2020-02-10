package etcd

import (
	"github.com/coreos/etcd/clientv3"
)

type Cfg struct {
	Muses struct{
		// reuse the clientv3,Config, it seems that we don't need to defined it again
		Etcd map[string] clientv3.Config `toml:"etcd"`
	} `toml:"muses"`
}
