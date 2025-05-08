package wrapper

import (
	"io"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// wrapper connection for gorilla package
type Connection struct {
	underlyingConnection *websocket.Conn

	// CloseTimeout amount of time to wait after sending a closure status
	CloseTimeout time.Duration
}

func New(conn *websocket.Conn) *Connection {
	return &Connection{
		underlyingConnection: conn,
		CloseTimeout:         time.Second * 30,
	}
}

// SendMessage implements websocket.Connection.
func (c *Connection) SendMessage(data []byte) error {
	return c.underlyingConnection.WriteMessage(websocket.TextMessage, data)
}

func (c *Connection) SendClose() error {
	return c.underlyingConnection.WriteControl(websocket.CloseMessage, nil, time.Now().Add(c.CloseTimeout))
}

func (c *Connection) ReceiveMessage() ([]byte, error) {
	_, data, err := c.underlyingConnection.ReadMessage()
	if err != nil && isClosedConnectionError(err) {
		return nil, io.EOF
	}
	return data, err
}

func (c *Connection) Close() error {
	err := c.underlyingConnection.Close()
	if err != nil && isClosedConnectionError(err) {
		return err
	}
	return nil
}

// check if close connection got error or not
func isClosedConnectionError(err error) bool {
	if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived, websocket.CloseGoingAway) {
		return true
	}
	return strings.Contains(err.Error(), "use of closed network connection")
}
