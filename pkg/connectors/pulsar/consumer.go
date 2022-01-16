// Package pulsar handles publishing and subscribing with Apache Pulsar
package pulsar

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/yngvark/gr-zombie/pkg/pubsub"
	"go.uber.org/zap"
)

type pulsarConsumer struct {
	log        *zap.SugaredLogger
	ctx        context.Context
	subscriber chan string
	client     pulsar.Client
	consumer   pulsar.Consumer
}

func (c *pulsarConsumer) SubscriberChannel() chan string {
	return c.subscriber
}

// ListenForMessages reads messages from Pulsar. This function blocks until the context provided on creation is done.
func (c *pulsarConsumer) ListenForMessages() error {
	select {
	case msg := <-c.consumer.Chan():
		msgString := string(msg.Payload())
		c.subscriber <- msgString
	case <-c.ctx.Done():
		return nil
	}

	return nil
}

func (c *pulsarConsumer) Close() error {
	c.log.Info("Closing pulsar consumer")

	c.consumer.Close()
	c.client.Close()

	return nil
}

const timeoutsDefault = 30 * time.Second

// NewConsumer returns a pulsar consumer
func NewConsumer(
	ctx context.Context,
	logger *zap.SugaredLogger,
	topic string,
	subscriber chan string,
) (pubsub.Consumer, error) {
	// Create client
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://localhost:36650",
		OperationTimeout:  timeoutsDefault,
		ConnectionTimeout: timeoutsDefault,
	})
	if err != nil {
		return nil, fmt.Errorf("could not instantiate Pulsar client: %w", err)
	}

	// Create consumer
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: "mysub2",
		Type:             pulsar.Exclusive,
	})
	if err != nil {
		return nil, fmt.Errorf("subscribing to client: %w", err)
	}

	c := &pulsarConsumer{
		log:        logger,
		ctx:        ctx,
		client:     client,
		consumer:   consumer,
		subscriber: subscriber,
	}

	return c, nil
}
