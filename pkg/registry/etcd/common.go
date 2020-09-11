package etcd

import (
	"github.com/mygomod/muses/pkg/common"
)

type Cfg struct {
	Muses struct {
		Etcd map[string]CallerCfg `toml:"etcd"`
	} `toml:"muses"`
}

// CallerCfg is similar with clientv3.Config
type CallerCfg struct {
	// Endpoints is a list of URLs.
	Endpoints []string

	// AutoSyncInterval is the interval to update endpoints with its latest members.
	// 0 disables auto-sync. By default auto-sync is disabled.
	AutoSyncInterval common.Duration

	// DialTimeout is the timeout for failing to establish a connection.
	DialTimeout common.Duration

	// DialKeepAliveTime is the time after which client pings the server to see if
	// transport is alive.
	DialKeepAliveTime common.Duration

	// DialKeepAliveTimeout is the time that the client waits for a response for the
	// keep-alive probe. If the response is not received in this time, the connection is closed.
	DialKeepAliveTimeout common.Duration

	// MaxCallSendMsgSize is the client-side request send limit in bytes.
	// If 0, it defaults to 2.0 MiB (2 * 1024 * 1024).
	// Make sure that "MaxCallSendMsgSize" < server-side default send/recv limit.
	// ("--max-request-bytes" flag to etcd or "embed.Config.MaxRequestBytes").
	MaxCallSendMsgSize int

	// MaxCallRecvMsgSize is the client-side response receive limit.
	// If 0, it defaults to "math.MaxInt32", because range response can
	// easily exceed request send limits.
	// Make sure that "MaxCallRecvMsgSize" >= server-side default send/recv limit.
	// ("--max-request-bytes" flag to etcd or "embed.Config.MaxRequestBytes").
	MaxCallRecvMsgSize int

	// Username is a user name for authentication.
	Username string

	// Password is a password for authentication.
	Password string

	// RejectOldCluster when set will refuse to create a client against an outdated cluster.
	RejectOldCluster bool

	// PermitWithoutStream when set will allow client to send keepalive pings to server without any active streams(RPCs).
	PermitWithoutStream bool
}
