package handler

import (
	"net/http"

	"github.com/boxuanduan/gochat/config"
	"github.com/boxuanduan/gochat/internal/gateway"
	"github.com/boxuanduan/gochat/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	hub      *gateway.Hub
	upgrader websocket.Upgrader
	jwtCfg   config.JWTConfig
}

func NewWSHandler(hub *gateway.Hub, jwtCfg config.JWTConfig) *WSHandler {
	return &WSHandler{
		hub: hub,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源，实际应用中应根据需求调整
			},
		},
		jwtCfg: jwtCfg,
	}
}

func (h *WSHandler) ServeWS(c *gin.Context) {
	// Get jwt token from query parameters
	token := c.Query("token")

	// Verify the token
	userID, err := auth.ParseToken(token, h.jwtCfg.Secret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	// Upgrade the HTTP connection to a WebSocket connection (No need to use c.Json here if error occurs, just return)
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// Create a new client and register it with the hub
	h.hub.HandleConn(conn, userID)
}
