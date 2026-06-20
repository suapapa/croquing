package ws

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait         = 60 * time.Second
	pingPeriod       = (pongWait * 9) / 10
	maxMessageSize   = 512
	clientSendBuffer = 256
)

// Client represents a single WebSocket connection in a lobby.
type Client struct {
	sync    *SnapshotSync
	lobbyID string
	conn    *websocket.Conn
	send    chan []byte
}

func newClient(sync *SnapshotSync, lobbyID string, conn *websocket.Conn) *Client {
	return &Client{
		sync:    sync,
		lobbyID: lobbyID,
		conn:    conn,
		send:    make(chan []byte, clientSendBuffer),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.sync.UnregisterClient(context.Background(), c.lobbyID, c)
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
