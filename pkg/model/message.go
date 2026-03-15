package model

import "time"

// Message 消息表（核心表）
type Message struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ConversationID string    `json:"conversation_id" gorm:"type:varchar(64);not null;index:idx_conv_seq"` // "p:{min}:{max}" 或 "g:{group_id}"
	SenderID       int64     `json:"sender_id" gorm:"not null"`
	MsgType        int8      `json:"msg_type" gorm:"default:0"` // 0=text, 1=image
	Content        string    `json:"content" gorm:"type:text;not null"`
	Seq            int64     `json:"seq" gorm:"not null;index:idx_conv_seq"` // 会话内递增序号
	CreatedAt      time.Time `json:"created_at"`
}

// UserConversation 用户会话状态（Timeline 游标）
type UserConversation struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         int64     `json:"user_id" gorm:"not null;uniqueIndex:uk_uc"`
	ConversationID string    `json:"conversation_id" gorm:"type:varchar(64);not null;uniqueIndex:uk_uc"`
	LastAckSeq     int64     `json:"last_ack_seq" gorm:"default:0"` // 已读位置
	UpdatedAt      time.Time `json:"updated_at"`
}

// Outbox 可靠投递表
type Outbox struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	MessageID  int64     `json:"message_id" gorm:"not null"`
	Status     int8      `json:"status" gorm:"default:0;index:idx_pending"` // 0=pending, 1=sent, 2=failed
	RetryCount int8      `json:"retry_count" gorm:"default:0"`
	CreatedAt  time.Time `json:"created_at" gorm:"index:idx_pending"`
}
