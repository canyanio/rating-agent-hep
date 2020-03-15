package rabbitmq

import (
	"context"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// ClientInterface is the interface for RabbitMQ client objects
type ClientInterface interface {
	GetMessageBusURI() string
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Publish(ctx context.Context, body string) error
}

// Client is the RabbitMQ client object
type Client struct {
	connection    *amqp.Connection
	messageBusURI string
}

// NewClient initializes a new HEP processor
func NewClient(messageBusURI string) *Client {
	return &Client{
		messageBusURI: messageBusURI,
	}
}

// GetMessageBusURI returns the message bus URI used by the client
func (c *Client) GetMessageBusURI() string {
	return c.messageBusURI
}

// Connect connects the client to the message bus
func (c *Client) Connect(ctx context.Context) error {
	connection, err := amqp.Dial(c.messageBusURI)
	if err != nil {
		return errors.Wrap(err, "unable to connect to the message bus")
	}
	c.connection = connection
	return nil
}

// Close closes the client to the message bus
func (c *Client) Close(ctx context.Context) error {
	err := c.connection.Close()
	if err != nil {
		return errors.Wrap(err, "unable to disconnect to the message bus")
	}
	return nil
}

// Publish publishes a message in the message bus
func (c *Client) Publish(ctx context.Context, body string) error {
	return nil
}
