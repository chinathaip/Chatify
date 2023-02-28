package service

import (
	"gorm.io/gorm"
)

type ChatService interface {
	GetAllChat() ([]Chat, error)
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
