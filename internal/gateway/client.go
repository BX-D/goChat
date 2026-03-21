package gateway

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub    *Hub
	userID int64
	conn   *websocket.Conn
	send   chan []byte
}

type WSMessage struct {
	From     int64  `json:"from"`
	To       int64  `json:"to"`
	Content  string `json:"content"`
	ChatType string `json:"chat_type"` // "private" or "group"
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

		// Save the message to the database

		if wsMessage.ChatType == "group" {
			convID := fmt.Sprintf("g:%d", wsMessage.To)
			// Store the message to the database with conversation ID
			_, err = c.hub.msgSvc.SendMessage(wsMessage.From, convID, wsMessage.Content)
			if err != nil {
				continue
			}
			// Get group members and push the message to each member
			members, err := c.hub.groupSvc.GetMemberIDs(wsMessage.To)
			if err != nil {
				continue
			}
			for _, memberID := range members {
				if memberID != c.userID { // Don't send the message back to the sender
					c.hub.Push(memberID, data)
				}
			}

		} else {
			convID := generateConversationID(wsMessage.From, wsMessage.To)
			_, err = c.hub.msgSvc.SendMessage(wsMessage.From, convID, wsMessage.Content)
			if err != nil {
				continue
			}
			// Push the message to the receiver
			c.hub.Push(wsMessage.To, data)
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

func generateConversationID(senderID int64, receiverID int64) string {
	conversationID := ""
	if senderID < receiverID {
		conversationID = fmt.Sprintf("p:%d:%d", senderID, receiverID)
	} else {
		conversationID = fmt.Sprintf("p:%d:%d", receiverID, senderID)
	}

	return conversationID
}
