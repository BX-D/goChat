package model

import "time"

// User 用户表
type User struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"type:varchar(64);uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"type:varchar(256);not null"` // json:"-" 表示序列化时隐藏密码
	Nickname  string    `json:"nickname" gorm:"type:varchar(64);default:''"`
	Avatar    string    `json:"avatar" gorm:"type:varchar(256);default:''"`
	CreatedAt time.Time `json:"created_at"`
}
