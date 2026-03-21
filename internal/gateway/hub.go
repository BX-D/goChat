package gateway

import (
	"github.com/boxuanduan/gochat/internal/chat"
	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[int64]*Client
	register   chan *Client
	unregister chan *Client
	msgSvc     *chat.MessageService
	groupSvc   *chat.GroupService
}

func NewHub(msgSvc *chat.MessageService, groupSvc *chat.GroupService) *Hub {
	return &Hub{
		clients:    make(map[int64]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		msgSvc:     msgSvc,
		groupSvc:   groupSvc,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.userID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
		}
	}
}

func (h *Hub) HandleConn(conn *websocket.Conn, userID int64) {
	// Create a new client
	client := &Client{
		hub:    h,
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, 256),
	}
	// Register the client with the hub
	h.register <- client

	// Start the client's read and write pumps
	go client.readPump()
	go client.writePump()
}

func (h *Hub) Push(userId int64, msg []byte) {
	if client, ok := h.clients[userId]; ok {
		client.send <- msg
	}
}
