package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/mendersoftware/go-lib-micro/log"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Client-specific constants
const (
	ExchangeName              = ""
	QueueNameBeginTransaction = "begin_transaction"
	QueueNameEndTransaction   = "end_transaction"
)

// ClientInterface is the interface for RabbitMQ client objects
type ClientInterface interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Publish(ctx context.Context, routingKey string, req interface{}) error
}

// Client is the RabbitMQ client object
type Client struct {
	connection    *amqp.Connection
	channel       *amqp.Channel
	messageBusURI string
}

// NewClient initializes a new HEP processor
func NewClient(messageBusURI string) *Client {
	return &Client{
		messageBusURI: messageBusURI,
	}
}

// Connect connects the client to the message bus
func (c *Client) Connect(ctx context.Context) error {
	l := log.FromContext(ctx)

	l.Infof("Connecting to message bus: %s", c.messageBusURI)
	connection, err := amqp.Dial(c.messageBusURI)
	if err != nil {
		return errors.Wrap(err, "unable to connect to the message bus")
	}
	c.connection = connection

	l.Info("Creating the message bus channel")
	channel, err := connection.Channel()
	if err != nil {
		return errors.Wrap(err, "unable to create the channel")
	}
	c.channel = channel

	for _, queue := range []string{QueueNameBeginTransaction, QueueNameEndTransaction} {
		l.Infof("Declaring the message bus queue: %s", queue)
		if _, err := channel.QueueDeclare(
			queue, // name
			false, // durable
			true,  // delete when unused
			false, // exclusive
			false, // no-wait
			amqp.Table{ // arguments
				"x-dead-letter-exchange": "rpc.dlx",
			},
		); err != nil {
			return errors.Wrapf(err, "unable to declare the queue: %s", queue)
		}
	}
	return nil
}

// Close closes the client to the message bus
func (c *Client) Close(ctx context.Context) error {
	if err := c.channel.Close(); err != nil {
		return errors.Wrap(err, "unable to close the channel")
	}
	if err := c.connection.Close(); err != nil {
		return errors.Wrap(err, "unable to disconnect to the message bus")
	}
	return nil
}

// Publish publishes a message in the message bus
func (c *Client) Publish(ctx context.Context, routingKey string, req interface{}) error {
	body, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "unable to marshal request to JSON")
	}
	err = c.channel.Publish(
		ExchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	if err != nil {
		return errors.Wrap(err, "unable to publish the request to the message bus")
	}
	return nil
}
