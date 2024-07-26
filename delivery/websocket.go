package delivery

import (
	"github.com/gorilla/websocket"
)

type websocketConn struct {
	conn *websocket.Conn
}

func NewWebsocketConn(conn *websocket.Conn) Delivery {
	return &websocketConn{conn: conn}
}

func (c *websocketConn) ReadMessage() ([]byte, error) {
	_, conn, err := c.conn.ReadMessage()
	return conn, err
}

func (c *websocketConn) WriteMessage(data []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *websocketConn) Close() error {
	return c.conn.Close()
}
