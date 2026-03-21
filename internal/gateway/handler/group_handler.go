package handler

import (
	"strconv"

	"github.com/boxuanduan/gochat/internal/chat"
	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	svc *chat.GroupService
}

func NewGroupHandler(svc *chat.GroupService) *GroupHandler {
	return &GroupHandler{svc: svc}
}

func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req struct {
		GroupName string `json:"group_name" binding:"required"`
		OwnerID   int64  `json:"owner_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	group, err := h.svc.CreateGroup(req.OwnerID, req.GroupName)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, group)
}

func (h *GroupHandler) JoinGroup(c *gin.Context) {
	var req struct {
		GroupID int64 `json:"group_id" binding:"required"`
		UserID  int64 `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.JoinGroup(req.GroupID, req.UserID)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "joined group successfully"})
}

func (h *GroupHandler) GetMembers(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseInt(groupIDStr, 0, 64)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	memberIDs, err := h.svc.GetMemberIDs(groupID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"member_ids": memberIDs})
}
