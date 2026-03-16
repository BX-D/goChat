package handler

import (
	"strconv"

	"github.com/boxuanduan/gochat/internal/chat"
	"github.com/gin-gonic/gin"
)

type MsgHandler struct {
	msgSvc *chat.MessageService
}

func NewMsgHandler(msgSvc *chat.MessageService) *MsgHandler {
	return &MsgHandler{msgSvc: msgSvc}
}

func (h *MsgHandler) PullMessages(c *gin.Context) {
	conversationID := c.Query("conversation_id")
	afterSeqStr := c.Query("after_seq")
	limitStr := c.Query("limit")

	afterSeq, err := strconv.ParseInt(afterSeqStr, 0, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid after_seq"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid limit"})
		return
	}

	msgs, err := h.msgSvc.PullMessages(conversationID, afterSeq, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, msgs)
}
