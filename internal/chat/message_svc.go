package chat

import (
	"context"
	"fmt"

	"github.com/boxuanduan/gochat/internal/repo/mysql"
	"github.com/boxuanduan/gochat/internal/repo/redis"
	"github.com/boxuanduan/gochat/pkg/model"
)

type MessageService struct {
	msgRepo *mysql.MessageRepo
	seqRepo *redis.SeqRepo
}

func NewMessageService(msgRepo *mysql.MessageRepo, seqRepo *redis.SeqRepo) *MessageService {
	return &MessageService{msgRepo: msgRepo, seqRepo: seqRepo}
}

// SendMessage: Generate conversationID -> generate seq -> save to db
func (s *MessageService) SendMessage(senderID, receiverID int64, content string) (*model.Message, error) {
	conversationID := generateConversationID(senderID, receiverID)
	seq, err := s.seqRepo.NextSeq(context.Background(), conversationID)
	if err != nil {
		return nil, fmt.Errorf("get max seq: %w", err)
	}

	msg := &model.Message{
		ConversationID: conversationID,
		Seq:            seq,
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

func (s *MessageService) AckMessage(userID int64, conversationID string, seq int64) error {
	// Call the repository method to update the acked seq for the user and conversation
	err := s.msgRepo.AckSeq(userID, conversationID, seq)
	return err
}
