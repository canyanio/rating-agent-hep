package state

import (
	"context"
	"flag"
	"testing"

	"github.com/mendersoftware/go-lib-micro/config"
	"github.com/stretchr/testify/assert"

	dconfig "github.com/canyanio/rating-agent-hep/config"
)

func TestNewRedisManager(t *testing.T) {
	flag.Parse()
	if testing.Short() {
		t.Skip()
	}

	redisAddress := config.Config.GetString(dconfig.SettingRedisAddress)
	redisPassword := config.Config.GetString(dconfig.SettingRedisPassword)
	redisDb := config.Config.GetInt(dconfig.SettingRedisDb)

	mgr := NewRedisManager(redisAddress, redisPassword, redisDb)
	assert.NotNil(t, mgr)

	ctx := context.Background()
	err := mgr.Connect(ctx)
	assert.Nil(t, err)

	err = mgr.Close(ctx)
	assert.Nil(t, err)
}

func TestNewRedisManagerConnectionFailed(t *testing.T) {
	flag.Parse()
	if testing.Short() {
		t.Skip()
	}

	mgr := NewRedisManager("localhost:1234", "", 0)
	assert.NotNil(t, mgr)

	ctx := context.Background()
	err := mgr.Connect(ctx)
	assert.NotNil(t, err)

	err = mgr.Close(ctx)
	assert.NotNil(t, err)
}
