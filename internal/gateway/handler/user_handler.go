package handler

import (
	"github.com/boxuanduan/gochat/internal/chat"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *chat.UserService
}

func NewUserHandler(svc *chat.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Register handles user registration requests.
func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.Register(req.Username, req.Password, req.Nickname)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := h.svc.Login(req.Username, req.Password)

	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token})
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/register", h.Register)
	r.POST("/api/login", h.Login)
}
