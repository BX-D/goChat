package chat

import (
	"fmt"

	"github.com/boxuanduan/gochat/internal/repo/mysql"
	"github.com/boxuanduan/gochat/pkg/model"
)

type MessageService struct {
	msgRepo *mysql.MessageRepo
}

func NewMessageService(msgRepo *mysql.MessageRepo) *MessageService {
	return &MessageService{msgRepo: msgRepo}
}

// SendMessage: Generate conversationID -> generate seq -> save to db
func (s *MessageService) SendMessage(senderID, receiverID int64, content string) (*model.Message, error) {
	conversationID := generateConversationID(senderID, receiverID)
	seq, err := s.msgRepo.GetMaxSeq(conversationID)
	if err != nil {
		return nil, fmt.Errorf("get max seq: %w", err)
	}

	nextSeq := seq + 1

	msg := &model.Message{
		ConversationID: conversationID,
		Seq:            nextSeq,
		SenderID:       senderID,
		Content:        content,
	}

	if err = s.msgRepo.Create(msg); err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}
	return msg, nil
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

func (s *MessageService) PullMessages(conversationID string, afterSeq int64, limit int) ([]*model.Message, error) {
	conversation, err := s.msgRepo.ListByConversation(conversationID, afterSeq, limit)
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}
	return conversation, nil
}
