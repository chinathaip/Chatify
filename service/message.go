package service

import (
	"time"

	"gorm.io/gorm"
)

type MessageService interface {
	GetMessagesInChat(int) ([]Message, error)
	StoreNewMessage(*Message) error
}

type MessageModel struct {
	DB *gorm.DB
}

func (m *MessageModel) GetMessagesInChat(chatID int) ([]Message, error) {
	var messages []Message
	err := m.DB.Preload("Sender").Where("chat_id = ?", chatID).Order("message_id ASC").Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (m *MessageModel) StoreNewMessage(msg *Message) error {
	msg.CreatedAt = time.Now()
	if err := m.DB.Create(msg).Where("chat_id=?", msg.ChatID).Error; err != nil {
		return err
	}

	return nil
}
