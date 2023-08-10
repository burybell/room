package room

import "github.com/gorilla/websocket"

var _ IConn = (*ConnWS)(nil)

// ConnWS Websocket connection
type ConnWS struct {
	conn *websocket.Conn
}

func NewConnWS(conn *websocket.Conn) *ConnWS {
	return &ConnWS{conn: conn}
}

func (r *ConnWS) Heartbeat() error {
	return r.conn.WriteMessage(websocket.PingMessage, []byte("heartbeat"))
}

func (r *ConnWS) PushMessage(message IMessage) error {
	bs, err := message.Bytes()
	if err != nil {
		return err
	}
	return r.conn.WriteMessage(websocket.TextMessage, bs)
}

func (r *ConnWS) Close() error {
	return r.conn.Close()
}
