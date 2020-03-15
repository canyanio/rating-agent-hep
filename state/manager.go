package state

import (
	"context"
)

// ManagerInterface describes a state manager object
type ManagerInterface interface {
	Connect(context context.Context) error
	Close(context context.Context) error
}
