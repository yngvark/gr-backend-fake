package httphandler

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/yngvark/gr-zombie/pkg/connectors"
	"github.com/yngvark/gr-zombie/pkg/pubsub/broadcast"
	"go.uber.org/zap"
	"net/http"
)

// New returns a HTTP handler that handles incoming websocket connections
// context is used to disconnect clients when the caller decides it's time to stop.
// subscriber is used as a channel to forward incoming messages from the websocket to
func New(
	ctx context.Context,
	logger *zap.SugaredLogger,
	allowedCorsOrigins map[string]bool,
	onConnect connectors.OnConnect,
	subscriber chan string,
	broadcaster *broadcast.Broadcaster,
) func(writer http.ResponseWriter, request *http.Request) {
	upgrader := &websocket.Upgrader{
		CheckOrigin:       createWebsocketCheckOriginFn(logger, allowedCorsOrigins),
		EnableCompression: true,
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		connection, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			logger.Error("could not upgrade:", err)
			return
		}

		logger.Info("Client connected!")

		h := NewConnectedHandler(ctx, logger, connection, subscriber, broadcaster)

		closeConnectionChannel := make(chan bool)
		messagesToClientChannel := make(chan string)

		go h.readIncomingMessages(closeConnectionChannel)
		go h.forwardMessagesToClient(closeConnectionChannel, messagesToClientChannel)

		err = onConnect(messagesToClientChannel)
		if err != nil {
			logger.Error("on connect:", err)
			return
		}

		h.closeConnectionWhenDone(closeConnectionChannel)
	}
}

type webfn func(r *http.Request) bool

func createWebsocketCheckOriginFn(logger *zap.SugaredLogger, allowedOrigins map[string]bool) webfn {
	return func(r *http.Request) bool {
		origin, ok := r.Header["Origin"]
		if !ok {
			return false
		}

		if len(origin) > 0 {
			_, ok := allowedOrigins[origin[0]]
			logger.Infof("Checking origin %s. Result: %t\n", origin[0], ok)

			return ok
		}

		return true
	}
}
