package models

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	AdminId   uint      `json:"adminId" gorm:"not null"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	CreatedAt time.Time `json:"createAt"`
}

type UserMessage struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	UserId    uint           `json:"userId" gorm:"not null;index"`
	MessageId uint           `json:"messageId" gorm:"not null;index"`
	Message   Message        `json:"message" gorm:"foreignKey:MessageId;references:ID"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"` // 用户删除消息时写入
}
