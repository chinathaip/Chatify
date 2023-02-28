package service

import (
	"time"

	"gorm.io/gorm"
)

type ChatService interface {
	GetAllChat() ([]Chat, error)
	CreateNewChat(*Chat) error
}

type ChatModel struct {
	DB *gorm.DB
}

func (m *ChatModel) GetAllChat() ([]Chat, error) {
	var chats []Chat
	if err := m.DB.Find(&chats).Error; err != nil {
		return nil, err
	}

	return chats, nil
}

func (m *ChatModel) CreateNewChat(chat *Chat) error {
	chat.CreatedAt = time.Now()
	if err := m.DB.Create(chat).Error; err != nil {
		return err
	}

	return nil
}
