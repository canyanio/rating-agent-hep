package state

import (
	"context"
)

// MemoryManager is the Memory state manager
type MemoryManager struct {
}

// NewMemoryManager returns a new Memory Manager objecft
func NewMemoryManager() *MemoryManager {
	return &MemoryManager{}
}

// Connect connects to the memory
func (m *MemoryManager) Connect(context context.Context) error {
	return nil
}

// Close disconnects from the memory
func (m *MemoryManager) Close(context context.Context) error {
	return nil
}
