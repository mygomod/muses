package registry

import (
	"context"

	"github.com/coreos/etcd/clientv3"

	"github.com/goecology/muses/pkg/app"
	"github.com/goecology/muses/pkg/logger"
	"github.com/goecology/muses/pkg/registry/etcd"
)

var (
	client            *clientv3.Client
	serverStatAddrKey string
)

func init() {
	client = etcd.Caller("default")
	if client == nil {
		logger.DefaultLogger().Error("Could not initialize the etcdv3 client")
	}
	serverStatAddrKey = "/muses/serverStateAddr/" + app.Config().Muses.App.AppKey()
}

func RegisterServerStatAddr(addr string) error {
	_, err := client.Put(context.Background(), serverStatAddrKey, addr)
	if err != nil {
		return err
	}
	return nil
}
