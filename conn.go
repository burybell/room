package room

// IConn Client long connection interface, already implemented with ws(ConnWS) and sse(ConnSSE)
type IConn interface {

	// Heartbeat Send heartbeat to client
	Heartbeat() error

	// PushMessage Send message to client
	PushMessage(message IMessage) error

	// Close Connection will be closed
	Close() error
}
