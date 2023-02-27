package db

import "gorm.io/gorm"

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
