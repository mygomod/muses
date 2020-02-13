package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallerCfg_AppKey(t *testing.T) {
	cfg := CallerCfg{
		Name: "MyApp",
	}

	assert.Equal(t, "MyApp", cfg.AppKey())
	cfg.Env = "Env"
	assert.Equal(t, "MyApp-Env", cfg.AppKey())
	cfg.Version = "Version"
	assert.Equal(t, "MyApp-Env-Version", cfg.AppKey())
}
