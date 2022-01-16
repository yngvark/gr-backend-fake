package pulsar

import (
	"context"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/yngvark/gr-zombie/pkg/pubsub"
	"go.uber.org/zap"
)

type pulsarPublisher struct {
	log      *zap.SugaredLogger
	ctx      context.Context
	cancelFn context.CancelFunc
	client   pulsar.Client
	producer pulsar.Producer
}

func (m *pulsarPublisher) SendMsg(msg string) error {
	_, err := m.producer.Send(m.ctx, &pulsar.ProducerMessage{
		Payload: []byte(msg),
	})
	if err != nil {
		m.cancelFn()
		return fmt.Errorf("sending message: %w", err)
	}

	return nil
}

func (m *pulsarPublisher) Close() error {
	m.log.Info("Closing pulsar publisher")
	m.producer.Close()
	m.client.Close()

	return nil
}

// NewPublisher returns a pulsar publisher
func NewPublisher(
	ctx context.Context,
	cancelFn context.CancelFunc,
	logger *zap.SugaredLogger,
	topic string,
) (pubsub.Publisher, error) {
	// Create client
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://localhost:36650",
		OperationTimeout:  timeoutsDefault,
		ConnectionTimeout: timeoutsDefault,
	})
	if err != nil {
		return nil, fmt.Errorf("could not instantiate Pulsar client: %w", err)
	}

	// Create producer
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create producer: %w", err)
	}

	p := &pulsarPublisher{
		log:      logger,
		ctx:      ctx,
		cancelFn: cancelFn,
		client:   client,
		producer: producer,
	}

	return p, nil
}
