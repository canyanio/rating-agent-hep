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

func TestRedisManagerGetSetDeleteInt(t *testing.T) {
	flag.Parse()
	if testing.Short() {
		t.Skip()
	}

	redisAddress := config.Config.GetString(dconfig.SettingRedisAddress)
	redisPassword := config.Config.GetString(dconfig.SettingRedisPassword)
	redisDb := config.Config.GetInt(dconfig.SettingRedisDb)

	mgr := NewRedisManager(redisAddress, redisPassword, redisDb)

	ctx := context.Background()
	mgr.Connect(ctx)
	defer mgr.Close(ctx)
	mgr.flushAll(ctx)

	var ret int
	err := mgr.Get(ctx, "key", &ret)
	assert.Nil(t, err)
	assert.Equal(t, 0, ret)

	mgr.Set(ctx, "key", 1, 0)

	ret = 0
	err = mgr.Get(ctx, "key", &ret)
	assert.Nil(t, err)
	assert.Equal(t, 1, ret)

	mgr.Delete(ctx, "key")

	ret = 0
	err = mgr.Get(ctx, "key", &ret)
	assert.Nil(t, err)
	assert.Equal(t, 0, ret)
}

func TestRedisManagerGetSetDeleteString(t *testing.T) {
	flag.Parse()
	if testing.Short() {
		t.Skip()
	}

	redisAddress := config.Config.GetString(dconfig.SettingRedisAddress)
	redisPassword := config.Config.GetString(dconfig.SettingRedisPassword)
	redisDb := config.Config.GetInt(dconfig.SettingRedisDb)

	mgr := NewRedisManager(redisAddress, redisPassword, redisDb)

	ctx := context.Background()
	mgr.Connect(ctx)
	defer mgr.Close(ctx)
	mgr.flushAll(ctx)

	var ret string
	err := mgr.Get(ctx, "key", &ret)
	assert.Nil(t, err)
	assert.Equal(t, "", ret)

	mgr.Set(ctx, "key", "TEST", 0)

	err = mgr.Get(ctx, "key", &ret)
	assert.Nil(t, err)
	assert.Equal(t, "TEST", ret)

	mgr.Delete(ctx, "key")

	ret = ""
	err = mgr.Get(ctx, "key", &ret)
	assert.Nil(t, err)
	assert.Equal(t, "", ret)
}
