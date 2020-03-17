package state

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoryManager(t *testing.T) {
	mgr := NewMemoryManager()
	assert.NotNil(t, mgr)

	ctx := context.Background()
	err := mgr.Connect(ctx)
	assert.Nil(t, err)

	err = mgr.Close(ctx)
	assert.Nil(t, err)
}

func TestMemoryManagerGetSetDeleteInt(t *testing.T) {
	mgr := NewMemoryManager()
	ctx := context.Background()

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

func TestMemoryManagerGetSetDeleteString(t *testing.T) {
	mgr := NewMemoryManager()
	ctx := context.Background()

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
