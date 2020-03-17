package state

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

// MemoryManager is the Memory state manager
type MemoryManager struct {
	store map[string]interface{}
}

// NewMemoryManager returns a new Memory Manager objecft
func NewMemoryManager() *MemoryManager {
	store := make(map[string]interface{})
	return &MemoryManager{
		store: store,
	}
}

// Connect connects to the memory
func (m *MemoryManager) Connect(context context.Context) error {
	return nil
}

// Close disconnects from the memory
func (m *MemoryManager) Close(context context.Context) error {
	return nil
}

// Set updates the data associated with a key
func (m *MemoryManager) Set(context context.Context, key string, data interface{}, ttl int) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "unable to marshal request to JSON")
	}
	m.store[key] = dataJSON
	return nil
}

// Get retrives the data associated with a key
func (m *MemoryManager) Get(context context.Context, key string, destination interface{}) error {
	dataJSON := m.store[key]
	if dataJSON == nil {
		return nil
	}
	err := json.Unmarshal(dataJSON.([]byte), destination)
	if err != nil {
		return errors.Wrap(err, "unable to marshal request to JSON")
	}
	return nil
}

// Delete deletes a key and its associated data
func (m *MemoryManager) Delete(context context.Context, key string) error {
	delete(m.store, key)
	return nil
}
