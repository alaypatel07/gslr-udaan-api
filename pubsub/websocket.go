package pubsub

import "github.com/gorilla/websocket"

type WsConn struct {
	*websocket.Conn
}

func NewWsConn(c *websocket.Conn) *WsConn {
	return &WsConn{c}
}

func (w WsConn) Consume(data interface{})  {
	w.WriteJSON(data)
}