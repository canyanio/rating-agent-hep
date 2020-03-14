package rabbitmq

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// ClientInterface is the interface for RabbitMQ client objects
type ClientInterface interface {
	Publish(body string) error
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

// Connect connects the client to the message bus
func (c *Client) Connect() error {
	connection, err := amqp.Dial(c.messageBusURI)
	if err != nil {
		return errors.Wrap(err, "unable to connect to the message bus")
	}
	c.connection = connection
	return nil
}

// Close closes the client to the message bus
func (c *Client) Close() error {
	err := c.connection.Close()
	if err != nil {
		return errors.Wrap(err, "unable to disconnect to the message bus")
	}
	return nil
}
