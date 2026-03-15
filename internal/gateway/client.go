package gateway

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub    *Hub
	userID int64
	conn   *websocket.Conn
	send   chan []byte
}

type WSMessage struct {
	From    int64  `json:"from"`
	To      int64  `json:"to"`
	Content string `json:"content"`
}

// readPump 负责从 WebSocket 连接读取消息并处理
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		// Parse the message and handle it
		wsMessage := &WSMessage{}
		err = json.Unmarshal(msg, wsMessage)
		if err != nil {
			continue
		}
		// Set the sender's user ID
		wsMessage.From = c.userID

		data, err := json.Marshal(wsMessage)

		if err != nil {
			continue
		}

		// Push the message to the recipient's client
		c.hub.Push(wsMessage.To, data)
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
