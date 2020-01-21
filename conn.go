package socketgo

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Conn struct {
	*websocket.Conn
}

func (c *Conn) Send(event string, payload interface{}) error {
	encoded, err := json.Marshal(NewMessage(event, payload))
	if err != nil {
		return err
	}

	return c.WriteMessage(websocket.TextMessage, encoded)
}

func (c *Conn) SendBytes(payload []byte) error {
	return c.WriteMessage(websocket.TextMessage, payload)
}

func NewConn(conn *websocket.Conn) *Conn {
	return &Conn{conn}
}
