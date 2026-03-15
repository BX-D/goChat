package model

import "time"

// Friend 好友关系表
type Friend struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64     `json:"user_id" gorm:"not null;uniqueIndex:uk_pair"`
	FriendID  int64     `json:"friend_id" gorm:"not null;uniqueIndex:uk_pair"`
	Status    int8      `json:"status" gorm:"default:0"` // 0=pending, 1=accepted
	CreatedAt time.Time `json:"created_at"`
}

// GroupInfo 群组表
type GroupInfo struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string    `json:"name" gorm:"type:varchar(128);not null"`
	OwnerID    int64     `json:"owner_id" gorm:"not null"`
	MaxMembers int       `json:"max_members" gorm:"default:500"`
	CreatedAt  time.Time `json:"created_at"`
}

// GroupMember 群成员表
type GroupMember struct {
	ID       int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupID  int64     `json:"group_id" gorm:"not null;uniqueIndex:uk_gm"`
	UserID   int64     `json:"user_id" gorm:"not null;uniqueIndex:uk_gm"`
	Role     int8      `json:"role" gorm:"default:0"` // 0=member, 1=admin, 2=owner
	JoinedAt time.Time `json:"joined_at"`
}
