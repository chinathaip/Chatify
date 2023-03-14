package service

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        int       `gorm:"column:chat_id" json:"id"`
	Name      string    `gorm:"column:chat_name" json:"chat_name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

type ChatParticipants struct {
	ID     int       `gorm:"column:chat_participant_id" json:"id"`
	ChatID int       `gorm:"foriegnKey:chat_id:" json:"chat_id"`
	UserID uuid.UUID `gorm:"foreignKey:user_id" json:"user_id"`
}

type Message struct {
	ID        int       `gorm:"column:message_id" json:"id"`
	Sender    User      `gorm:"foreignKey:sender_id" json:"sender"`
	SenderID  uuid.UUID `gorm:"column:sender_id" json:"-"`
	ChatID    int       `gorm:"column:chat_id" json:"chat_id"`
	Data      string    `gorm:"column:data" json:"data"`
	CreatedAt time.Time `gorm:"column:sent_at" json:"sent_at"`
}

type User struct {
	ID       uuid.UUID `gorm:"column:id" json:"id"`
	Username string    `gorm:"column:username" json:"username"`
}
