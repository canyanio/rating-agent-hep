package state

import (
	"context"
)

// ManagerInterface describes a state manager object
type ManagerInterface interface {
	Connect(context context.Context) error
	Close(context context.Context) error
	Set(context context.Context, key string, req interface{}, ttl int) error
	Get(context context.Context, key string, destination *interface{}) error
	Delete(context context.Context, key string) error
}
