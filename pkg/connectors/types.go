// Package connectors contains ways to communicate with different technologies, for instance websockets or
// Kafka.
package connectors

// Connector is used to connect to clients. Implementors can use websockets, pulsar, kafka, etc.
type Connector interface {
	// ListenForConnections listens for incoming connections. It may or may not block, see implementation comments.
	ListenForConnections(OnConnect) error
	// Close closes the Connector
	Close() error
}

// OnConnect is a function that is called when a client connects
type OnConnect func(messagesToClientChannel chan string) error
