package main

import (
	"context"
	"fmt"
	"github.com/yngvark/gr-zombie/pkg/connectors/kafka"
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
	context   context.Context
	logger    *zap.SugaredLogger
	publisher pubsub.Publisher
	consumer  pubsub.Consumer
}

type getEnv func(key string) string

func newGameOpts(ctx context.Context, cancelFn context.CancelFunc, getEnv getEnv) (*GameOpts, error) {
	logger, err := log2.New()
	if err != nil {
		return nil, fmt.Errorf("could not create logger: %w", err)
	}

	var publisher pubsub.Publisher

	var consumer pubsub.Consumer

	subscriber := make(chan string)

	switch {
	case getEnv("GAME_QUEUE_TYPE") == "kafka":
		publisher, consumer, err = pubSubForKafka(ctx, cancelFn, logger, subscriber)
		if err != nil {
			return nil, fmt.Errorf("creating pulsar connectors: %w", err)
		}
	default:
		publisher, consumer, err = pubSubForWebsocket(ctx, cancelFn, logger, subscriber)
		if err != nil {
			return nil, fmt.Errorf("creating websocket connectors: %w", err)
		}
	}

	return &GameOpts{
		context:   ctx,
		logger:    logger,
		publisher: publisher,
		consumer:  consumer,
	}, nil
}

const allowedCorsOriginsEnvVarKey = "ALLOWED_CORS_ORIGINS"

func pubSubForWebsocket(
	ctx context.Context,
	cancelFn context.CancelFunc,
	logger *zap.SugaredLogger,
	subscriber chan string,
) (pubsub.Publisher, pubsub.Consumer, error) {
	corsHelper := oslookup.NewCORSHelper(logger)

	allowedCorsOrigins, err := corsHelper.GetAllowedCorsOrigins(os.LookupEnv, allowedCorsOriginsEnvVarKey)
	if err != nil {
		return nil, nil, fmt.Errorf("getting allowed CORS origins: %w", err)
	}

	corsHelper.PrintAllowedCorsOrigins(allowedCorsOrigins)

	p, c := websocket.New(ctx, cancelFn, logger, subscriber, allowedCorsOrigins)

	return p, c, nil
}

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

func pubSubForKafka(
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
