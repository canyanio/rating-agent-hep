package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Client is the RabbitMQ client mock
type Client struct {
	mock.Mock
}

// Connect connects the client to the message bus
func (c *Client) Connect(ctx context.Context) error {
	ret := c.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// Close closes the client to the message bus
func (c *Client) Close(ctx context.Context) error {
	ret := c.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// Publish publishes a message in the message bus
func (c *Client) Publish(ctx context.Context, routingKey string, req interface{}) error {
	ret := c.Called(ctx, routingKey, req)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}
