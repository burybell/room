package room

import (
	"net/http"
)

var _ IConn = (*ConnSSE)(nil)

// ConnSSE Server-Sent Events connection
type ConnSSE struct {
	http.ResponseWriter
	http.Flusher
}

func NewConnSSE(responseWriter http.ResponseWriter, flusher http.Flusher) *ConnSSE {
	return &ConnSSE{ResponseWriter: responseWriter, Flusher: flusher}
}

func (r *ConnSSE) Heartbeat() error {
	return nil
}

func (r *ConnSSE) PushMessage(message IMessage) error {
	bs, err := message.Bytes()
	if err != nil {
		return err
	}
	_, err = r.ResponseWriter.Write([]byte("data: " + string(bs) + "\n\n"))
	if err != nil {
		return err
	}
	r.Flusher.Flush()
	return nil
}

func (r *ConnSSE) Close() error {
	r.Flusher.Flush()
	return nil
}
