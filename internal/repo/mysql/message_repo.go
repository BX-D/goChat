package mysql

import (
	"github.com/boxuanduan/gochat/pkg/model"
	"gorm.io/gorm"
)

type MessageRepo struct {
	db *gorm.DB
}

func NewMessageRepo(db *gorm.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) Create(msg *model.Message) error {
	tx := r.db.Create(msg)
	return tx.Error
}

func (r *MessageRepo) ListByConversation(conversationID string, afterSeq int64, limit int) ([]*model.Message, error) {
	msgs := make([]*model.Message, 0)
	err := r.db.Where("conversation_id = ? AND seq > ?", conversationID, afterSeq).
		Order("seq ASC").
		Limit(limit).
		Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (r *MessageRepo) GetMaxSeq(conversationID string) (int64, error) {
	var maxSeq int64

	err := r.db.Model(&model.Message{}).
		Where("conversation_id = ?", conversationID).
		Select("COALESCE(MAX(seq), 0)").
		Scan(&maxSeq).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil // No messages, return 0
		}
		return 0, err
	}

	return maxSeq, nil
}

func (r *MessageRepo) AckSeq(userID int64, conversationID string, seq int64) error {
	// Update the UserConversation record to set the acked seq for the user and conversation
	// If not exists, create a new record with the acked seq
	// If exists, update the acked seq if the new seq is greater than the existing one
	var record model.UserConversation
	err := r.db.Where("user_id = ? AND conversation_id = ?", userID, conversationID).
		Assign(model.UserConversation{
			UserID:         userID,
			ConversationID: conversationID,
			LastAckSeq:     seq,
		}).
		FirstOrCreate(&record).
		Error

	return err
}
