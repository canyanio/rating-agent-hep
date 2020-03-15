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
