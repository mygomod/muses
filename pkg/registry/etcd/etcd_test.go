package etcd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	assert.NotNil(t, Register())
}

func TestCallerStore_InitCfg(t *testing.T) {
	config := `
[muses.etcd.default]
    endpoints = ["endpoint1", "endpoint2", "endpoint3"]
    autoSyncInterval = "10s"
    dialTimeout = "9s"
    dialKeepAliveTime = "8s"
    dialKeepAliveTimeout = "7s"
    username = "my name"
    password = "my pwd"
    maxCallSendMsgSize = 20
    maxCallRecvMsgSize = 10
    rejectOldCluster = true
    permitWithoutStream = false

`
	store := Register()
	err := store.InitCfg([]byte(config))
	assert.Nil(t, err)
	cfg := store.(*callerStore).cfg
	assert.NotNil(t, cfg)
	assert.Equal(t, 1, len(cfg.Muses.Etcd))
	etcd, ok := cfg.Muses.Etcd["default"]
	assert.True(t, ok)
	assert.NotNil(t, etcd)
	assert.Equal(t, 3, len(etcd.Endpoints))
	assert.Equal(t, 10*time.Second, etcd.AutoSyncInterval.Duration)
	assert.Equal(t, 9*time.Second, etcd.DialTimeout.Duration)
	assert.Equal(t, 8*time.Second, etcd.DialKeepAliveTime.Duration)
	assert.Equal(t, 7*time.Second, etcd.DialKeepAliveTimeout.Duration)
	assert.Equal(t, "my name", etcd.Username)
	assert.Equal(t, "my pwd", etcd.Password)
	assert.Equal(t, 20, etcd.MaxCallSendMsgSize)
	assert.Equal(t, 10, etcd.MaxCallRecvMsgSize)
	assert.True(t, etcd.RejectOldCluster)
	assert.False(t, etcd.PermitWithoutStream)

	// test empty config

	err = store.InitCfg([]byte(""))
	assert.Nil(t, err)

	config = `
[muses.etcd.default]
`
	err = store.InitCfg([]byte(config))
	assert.Nil(t, err)

	// invalid config
	config = `
[muses.etcd.default]
    rejectOldCluster = aaa
`
	err = store.InitCfg([]byte(config))
	assert.NotNil(t, err)
}
