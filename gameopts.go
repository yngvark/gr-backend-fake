package main

import (
	"context"
	"fmt"
	"github.com/yngvark/gr-zombie/pkg/connectors"
	"github.com/yngvark/gr-zombie/pkg/pubsub/broadcast"
	"os"

	"github.com/yngvark/gr-zombie/pkg/connectors/websocket/oslookup"

	"github.com/yngvark/gr-zombie/pkg/connectors/pulsar"
	"github.com/yngvark/gr-zombie/pkg/connectors/websocket"
	"github.com/yngvark/gr-zombie/pkg/log2"
	"github.com/yngvark/gr-zombie/pkg/pubsub"
	"go.uber.org/zap"
)

// GameOpts contains various dependencies
type GameOpts struct {
	context     context.Context
	cancelFn    context.CancelFunc
	log         *zap.SugaredLogger
	subscriber  chan string
	broadcaster *broadcast.Broadcaster
	connector   connectors.Connector
}

type getEnv func(key string) string

//goland:noinspection GoUnusedParameter
func newGameOpts(ctx context.Context, cancelFn context.CancelFunc, getEnv getEnv) (*GameOpts, error) {
	log, err := log2.New()
	if err != nil {
		return nil, fmt.Errorf("could not create logger: %w", err)
	}

	broadcaster := broadcast.New(log)

	var connector connectors.Connector

	subscriber := make(chan string)

	switch {
	//case getEnv("GAME_QUEUE_TYPE") == "kafka":
	//	consumer, err = pubSubForKafka(ctx, cancelFn, logger, subscriber)
	//	if err != nil {
	//		return nil, fmt.Errorf("creating pulsar connectors: %w", err)
	//	}
	default:
		connector, err = newWebsocketConnector(ctx, log, subscriber, broadcaster)
		if err != nil {
			return nil, fmt.Errorf("creating websocket connectors: %w", err)
		}
	}

	return &GameOpts{
		context:     ctx,
		cancelFn:    cancelFn,
		log:         log,
		subscriber:  subscriber,
		broadcaster: broadcaster,
		connector:   connector,
	}, nil
}

const allowedCorsOriginsEnvVarKey = "ALLOWED_CORS_ORIGINS"

func newWebsocketConnector(
	ctx context.Context,
	logger *zap.SugaredLogger,
	subscriber chan string,
	broadcaster *broadcast.Broadcaster,
) (connectors.Connector, error) {
	corsHelper := oslookup.NewCORSHelper(logger)

	allowedCorsOrigins, err := corsHelper.GetAllowedCorsOrigins(os.LookupEnv, allowedCorsOriginsEnvVarKey)
	if err != nil {
		return nil, fmt.Errorf("getting allowed CORS origins: %w", err)
	}

	corsHelper.PrintAllowedCorsOrigins(allowedCorsOrigins)

	c := websocket.NewConnector(ctx, logger, subscriber, allowedCorsOrigins, broadcaster)

	return c, nil
}

//goland:noinspection GoUnusedFunction
func pubSubForPulsar(
	ctx context.Context,
	cancelFn context.CancelFunc,
	logger *zap.SugaredLogger,
	subscriber chan string,
) (pubsub.Publisher, pubsub.Consumer, error) {
	p, err := pulsar.NewPublisher(ctx, cancelFn, logger, "zombie")
	if err != nil {
		return nil, nil, fmt.Errorf("creating publisher: %w", err)
	}

	c, err := pulsar.NewConsumer(ctx, logger, "gameinit", subscriber)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create consumer: %w", err)
	}

	return p, c, nil
}

/*func pubSubForKafka(
	ctx context.Context,
	cancelFn context.CancelFunc,
	logger *zap.SugaredLogger,
	subscriber chan string,
) (pubsub.Publisher, pubsub.Consumer, error) {
	p, err := kafka.NewPublisher(ctx, cancelFn, logger, "zombie")
	if err != nil {
		return nil, nil, fmt.Errorf("creating publisher: %w", err)
	}

	c, err := kafka.NewConsumer(ctx, logger, "gameinit", subscriber)

	return p, c, nil
}
*/
