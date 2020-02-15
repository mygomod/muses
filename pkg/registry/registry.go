package registry

import (
	"context"
	"errors"

	"github.com/coreos/etcd/clientv3"
)

var (
	client            *clientv3.Client
	serverStatAddrKey string
)

// todo fix
//func init() {
//	client = etcd.Caller("default")
//	if client == nil {
//		logger.DefaultLogger().Error("Could not initialize the etcdv3 client")
//	}
//	serverStatAddrKey = "/muses/serverStateAddr/" + app.Config().Muses.App.AppKey()
//}

func RegisterServerStatAddr(addr string) error {
	if client == nil {
		return errors.New("etcd not register")
	}
	_, err := client.Put(context.Background(), serverStatAddrKey, addr)
	if err != nil {
		return err
	}
	return nil
}
