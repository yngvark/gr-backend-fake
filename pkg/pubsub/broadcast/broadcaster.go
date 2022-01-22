// Package broadcast knows how to broadcast messages to subscribers
package broadcast

import "go.uber.org/zap"

// Broadcaster is used for sending (broadcasting) messages to a number of subscribers
type Broadcaster struct {
	subscribers []chan<- string
	log         *zap.SugaredLogger
}

// AddSubscriber adds a Subscriber to its list of subscribers
func (b *Broadcaster) AddSubscriber(subscriber chan<- string) {
	b.subscribers = append(b.subscribers, subscriber)
}

// BroadCast sends a message to all Subscriber-s
func (b *Broadcaster) BroadCast(msg string) error {
	for _, subscriber := range b.subscribers {
		subscriber <- msg
	}

	return nil
}

// New returns a new Broadcaster
func New(logger *zap.SugaredLogger) *Broadcaster {
	return &Broadcaster{
		log:         logger,
		subscribers: make([]chan<- string, 0),
	}
}
