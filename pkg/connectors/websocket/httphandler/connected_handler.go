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

func (h *ConnectedHandler) readIncomingMessages(closeConnectionChannel chan<- bool) {
	for {
		h.log.Debug("Reading next message from client...")

		_, message, err := h.connection.ReadMessage()
		if err != nil {
			// Client disconnected
			h.log.Info("Client disconnected")

			closeConnectionChannel <- true

			closeError, ok := err.(*websocket.CloseError)
			if ok {
				h.log.Infof("Client disconnected OK. Code: %d", closeError.Code)
			} else {
				h.log.Errorf("Read error: %s", err.Error())
			}

			break
		}

		h.log.Infof("Sending received message to subscriber: %s", message)
		h.subscriber <- string(message)
	}

	h.log.Debug("Done reading incoming messages")
}

func (h *ConnectedHandler) closeConnectionWhenDone(closeConnectionChannel <-chan bool) {
	select {
	case <-closeConnectionChannel:
		h.log.Debug("ConnectedHandler.closeConnectionWhenDone")
	case <-h.ctx.Done():
		h.log.Debug("ConnectedHandler.ctx.Done")
	}

	h.log.Info("Closing connection to client")
	err := h.Close()

	if err != nil {
		h.log.Info("error when closing connection: %w", err)
	} else {
		h.log.Info("Connection closed successfully.")
	}
}

func (h *ConnectedHandler) forwardMessagesToClient(closeConnectionChannel chan bool, messagesToClientChannel chan string) {
	h.broadcaster.AddSubscriber(messagesToClientChannel)

	for {
		select {
		case msgToClient := <-messagesToClientChannel:
			err := h.sendMsgToConnection(msgToClient)
			if err != nil {
				h.log.Errorf("Could not send message to client. Stopping. Msg: %s", msgToClient)
				closeConnectionChannel <- true

				return
			}

		case <-closeConnectionChannel:
			h.log.Debug("closeConnectionChannel signal. Stopping broadcasting to client.")

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

// Close closes the handler
func (h *ConnectedHandler) Close() error {
	h.log.Info("Closing ConnectedHandler")

	if h.connection != nil {
		err := h.connection.Close()
		if err != nil {
			_, ok := err.(*net.OpError)
			if ok {
				return nil
			}

			return err
		}
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
