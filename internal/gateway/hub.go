package gateway

import "github.com/gorilla/websocket"

type Hub struct {
	clients    map[int64]*Client
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
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
