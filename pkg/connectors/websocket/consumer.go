package websocket

import (
	"context"
	"errors"
	"github.com/yngvark/gr-zombie/pkg/connectors"
	"github.com/yngvark/gr-zombie/pkg/pubsub/broadcast"
	"net/http"

	"github.com/yngvark/gr-zombie/pkg/connectors/websocket/httphandler"
	"go.uber.org/zap"
)

type connctionHandler struct {
	ctx        context.Context
	log        *zap.SugaredLogger
	subscriber chan string

	listening          bool
	broadcaster        *broadcast.Broadcaster
	allowedCorsOrigins map[string]bool
}

// ListenForConnections starts to receive messages which will be available by reading SubscriberChannel().
func (c *connctionHandler) ListenForConnections(onConnect connectors.OnConnect) error {
	if !c.listening {
		c.listening = true
	} else {
		return errors.New("already listening for messages. Can listen for messages only once")
	}

	http.HandleFunc(
		"/zombie",
		httphandler.New(c.ctx, c.log, c.allowedCorsOrigins, onConnect, c.subscriber, c.broadcaster),
	)

	return nil
}

func (c *connctionHandler) StopListening() error {
	c.log.Info("connctionHandler.StopListening")
	return nil
}

// NewConnector returns a new consumer for websockets
func NewConnector(
	ctx context.Context,
	logger *zap.SugaredLogger,
	subscriber chan string,
	allowedCorsOrigins map[string]bool,
	broadcaster *broadcast.Broadcaster,
) connectors.Connector {
	return &connctionHandler{
		ctx:                ctx,
		log:                logger,
		subscriber:         subscriber,
		broadcaster:        broadcaster,
		allowedCorsOrigins: allowedCorsOrigins,
	}
}
