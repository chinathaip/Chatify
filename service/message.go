package service

import (
	"time"

	"gorm.io/gorm"
)

type MessageService interface {
	GetMessagesInChat(int, int, int) ([]Message, error)
	StoreNewMessage(*Message) error
}

type MessageModel struct {
	DB *gorm.DB
}

func (m *MessageModel) GetMessagesInChat(chatID, pageNumber, pageSize int) ([]Message, error) {
	var messages []Message
	offset := (pageNumber - 1) * pageSize //skip the first n(pageSize) message, and start on the next one
	err := m.DB.Preload("Sender").Where("chat_id = ?", chatID).Offset(offset).Limit(pageSize).Order("message_id DESC").Find(&messages).Error
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
