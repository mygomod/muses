package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallerStore_InitCfg(t *testing.T) {
	store := Register()
	config := `
[muses.token.default]
	mode = "redis"

[muses.token.default.redis]
	mode = "redis"
    debug = true
    addr = "127.0.0.1:6379"
    network = "tcp"
    db = 0
    password = ""

`
	err := store.InitCfg([]byte(config))
	assert.Nil(t, err)
	err = store.InitCaller()
	assert.Nil(t, err)
}