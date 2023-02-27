package service

import (
	"time"
)

type Chat struct {
	ID        int       `gorm:"column:chat_id"`
	Name      int       `gorm:"column:chat_name"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

type Message struct {
	ID       int       `gorm:"column:message_id" json:"id"`
	SenderID int       `gorm:"column:sender_id" json:"sender_id"`
	Data     string    `gorm:"column:data" json:"data"`
	SentAt   time.Time `gorm:"column:sent_at" json:"sent_at"`
}
