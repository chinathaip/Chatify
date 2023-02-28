package service

import (
	"time"

	"gorm.io/gorm"
)

type MessageModel struct {
	DB *gorm.DB
}

func (m *MessageModel) GetMessagesInChat(chatID int) ([]Message, error) {
	var msg []Message
	result := m.DB.Where("chat_id=?", chatID).Find(&msg)
	if result.Error != nil {
		return nil, result.Error
	}

	return msg, nil
}

func (m *MessageModel) StoreNewMessage(msg *Message) error {
	msg.SentAt = time.Now()
	if err := m.DB.Create(msg).Where("chat_id=?", msg.ChatID).Error; err != nil {
		return err
	}

	return nil
}
