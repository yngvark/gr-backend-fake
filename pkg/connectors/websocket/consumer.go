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
	logger     *zap.SugaredLogger
	subscriber chan string

	httpHandler        *httphandler.ConnectedHandler
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
		httphandler.New(c.ctx, c.logger, c.allowedCorsOrigins, onConnect, c.subscriber, c.broadcaster),
	)

	//<-c.ctx.Done()

	return nil
}

// Close closes the connctionHandler
func (c *connctionHandler) Close() error {
	c.logger.Info("Closing connctionHandler")

	if c.httpHandler != nil {
		return c.httpHandler.Close()
	}

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
		logger:             logger,
		subscriber:         subscriber,
		broadcaster:        broadcaster,
		allowedCorsOrigins: allowedCorsOrigins,
	}
}
