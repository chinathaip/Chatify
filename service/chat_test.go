package service

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func seedChatDB(db *gorm.DB) {
	db.AutoMigrate(&Chat{}) //create necessary table in the temp db
	chat1 := Chat{ID: 1, Name: "Test Chat Room 1", CreatedAt: time.Now()}
	chat2 := Chat{ID: 2, Name: "Test Chat Room 2", CreatedAt: time.Now()}

	if err := db.Create(&chat1).Error; err != nil {
		log.Fatalln("cannot seed chat 1", err)
	}

	if err := db.Create(&chat2).Error; err != nil {
		log.Fatalln("cannot seed chat 2", err)
	}
}

func TestGetAllChat(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	chatModel := &ChatModel{DB: db}
	seedChatDB(db)
	expected := []Chat{
		{
			ID:   1,
			Name: "Test Chat Room 1",
		},
		{
			ID:   2,
			Name: "Test Chat Room 2",
		},
	}

	chats, err := chatModel.GetAllChat()

	assert.NoError(t, err)
	assert.Equal(t, expected[0].ID, chats[0].ID)
	assert.Equal(t, expected[1].Name, chats[1].Name)
}
