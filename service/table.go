package service

import (
	"time"
)

type Chat struct {
	ID        int       `gorm:"column:chat_id" json:"id"`
	Name      string    `gorm:"column:chat_name" json:"chat_name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

type Message struct {
	ID       int       `gorm:"column:message_id" json:"id"`
	SenderID string    `gorm:"column:sender_id" json:"sender_id"`
	ChatID   int       `gorm:"column:chat_id" json:"chat_id"` //not sent to client
	Data     string    `gorm:"column:data" json:"data"`
	SentAt   time.Time `gorm:"column:sent_at" json:"sent_at"`
}
