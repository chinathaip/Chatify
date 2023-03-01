package service

import (
	"log"
	"time"

	er "github.com/chinathaip/chatify/error"
	"gorm.io/gorm"
)

var chatErr = er.ChatError{}

type ChatService interface {
	GetAllChat() ([]Chat, error)
	CreateNewChat(*Chat) error
	DeleteChat(chatID int) error
	IsChatExist(chatName string) (int, bool)
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

func (m *ChatModel) DeleteChat(chatID int) error {
	if err := m.DB.Where("chat_id=?", chatID).Delete(&Chat{}).Error; err != nil {
		return err
	}

	return nil
}

func (m *ChatModel) IsChatExist(chatName string) (int, bool) {
	var chat Chat
	if err := m.DB.Where("chat_name=?", chatName).First(&chat).Error; err != nil {
		chatErr.Log(err)
		return 0, false
	}

	log.Printf("Got Chat : %v\n", chat)

	if chat.ID == 0 {
		return 0, false
	}

	return chat.ID, true
}
