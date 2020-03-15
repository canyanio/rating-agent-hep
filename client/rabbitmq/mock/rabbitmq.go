package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// Client is the RabbitMQ client mock
type Client struct {
	mock.Mock
}

// GetMessageBusURI returns the message bus URI used by the client
func (c *Client) GetMessageBusURI() string {
	ret := c.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}
	return r0
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
func (c *Client) Publish(ctx context.Context, body string) error {
	ret := c.Called(ctx, body)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(body)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}
