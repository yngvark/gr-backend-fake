// Package httphandler knows how to handle HTTP websocket connections
package httphandler

import (
	"context"
	"errors"
	"fmt"
	"github.com/yngvark/gr-zombie/pkg/pubsub/broadcast"
	"go.uber.org/zap"
	"net"

	"github.com/gorilla/websocket"
)

// ConnectedHandler knows how to handle a specific, connected HTTP websocket connection.
// It will be used when connection to a client has already been made.
type ConnectedHandler struct {
	log         *zap.SugaredLogger
	ctx         context.Context
	connection  *websocket.Conn
	subscriber  chan string
	broadcaster *broadcast.Broadcaster
}

func (h *ConnectedHandler) readIncomingMessages() {
	for {
		h.log.Debug("Reading next message from client...")

		_, message, err := h.connection.ReadMessage()
		if err != nil {
			// Client disconnected
			h.log.Info("Client disconnected")

			closeError, ok := err.(*websocket.CloseError)
			if ok {
				h.log.Debugf("Client disconnected OK. Code: %d", closeError.Code)
			} else {
				h.log.Errorf("Client disconnect read error: %s", err.Error())
			}

			return
		}

		h.log.Infof("Sending received message to subscriber: %s", message)
		h.subscriber <- string(message)
	}
}

func (h *ConnectedHandler) closeConnectionWhenDone(websocketReadStoppedChannel <-chan bool) {
	select {
	case <-websocketReadStoppedChannel:
		h.log.Debug("ConnectedHandler.closeConnectionWhenDone.websocketReadFailureChannel")
		return
	case <-h.ctx.Done():
		h.log.Debug("ConnectedHandler.closeConnectionWhenDone.ctx.Done")
	}

	h.log.Info("Closing connection to client")
	err := h.CloseIt()

	if err != nil {
		h.log.Error("ConnectedHandler.closeConnectionWhenDone: %w", err)
	} else {
		h.log.Info("ConnectedHandler.closeConnectionWhenDone success")
	}
}

func (h *ConnectedHandler) forwardMessagesToClient(messagesToClientChannel chan string) {
	h.broadcaster.AddSubscriber(messagesToClientChannel)

	for {
		select {
		case msgToClient := <-messagesToClientChannel:
			err := h.sendMsgToConnection(msgToClient)
			if err != nil {
				h.log.Info("Could not send message to client. Stopping handler for this connection.")

				err = h.CloseIt()
				if err != nil {
					h.log.Errorf("error closing: %w", err)
				}

				return
			}
		case <-h.ctx.Done():
			h.log.Debug("ConnectedHandler.forwardMessagesToClient.ctx.Done. Stopping broadcasting to client.")
			return
		}
	}
}

// sendMsgToConnection sends a message via the websocket
func (h *ConnectedHandler) sendMsgToConnection(msg string) error {
	if h.connection == nil {
		return errors.New("could not send message, not connected")
	}

	//h.log.Debugf("Sending msg: %s", msg)

	err := h.connection.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return fmt.Errorf("could not write message: %w", err)
	}

	return nil
}

// CloseIt closes the handler
func (h *ConnectedHandler) CloseIt() error {
	h.log.Info("ConnectedHandler.Close()")

	if h.connection != nil {
		err := h.connection.Close()
		if err != nil {
			_, ok := err.(*net.OpError)
			if ok {
				return nil
			}

			return err
		}

		h.log.Info("ConnectedHandler.Close() success.")
	}

	return nil
}

// NewConnectedHandler returns a new ConnectedHandler
func NewConnectedHandler(
	ctx context.Context,
	logger *zap.SugaredLogger,
	connection *websocket.Conn,
	subscriber chan string,
	broadcaster *broadcast.Broadcaster,
) *ConnectedHandler {
	handler := &ConnectedHandler{
		ctx:         ctx,
		log:         logger,
		connection:  connection,
		subscriber:  subscriber,
		broadcaster: broadcaster,
	}

	return handler
}
