package gateway

import "github.com/gorilla/websocket"

type Client struct {
	hub    *Hub
	userID int64
	conn   *websocket.Conn
	send   chan []byte
}

// readPump 负责从 WebSocket 连接读取消息并处理
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
