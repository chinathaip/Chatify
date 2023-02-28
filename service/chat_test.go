package service

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func clearChatDB(db *gorm.DB) {
	if err := db.Delete(&Chat{}, "chat_id>=1").Error; err != nil {
		log.Fatal(err)
	}
}

func TestGetAllChat(t *testing.T) {
	db, dbConn := setup()
	defer teardown(db, dbConn, clearChatDB)
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

func TestCreateNewChat(t *testing.T) {
	db, dbConn := setup()
	defer teardown(db, dbConn, clearChatDB)
	chatModel := &ChatModel{DB: db}
	seedChatDB(db)
	newChat := &Chat{Name: "Test Chat Room 3"}
	expected := []Chat{
		{
			ID:   1,
			Name: "Test Chat Room 1",
		},
		{
			ID:   2,
			Name: "Test Chat Room 2",
		},
		{
			ID:   3,
			Name: "Test Chat Room 3",
		},
	}

	err := chatModel.CreateNewChat(newChat)
	result, _ := chatModel.GetAllChat()

	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(result))
	assert.Equal(t, expected[2].Name, result[2].Name)
}

func TestDeleteChat(t *testing.T) {
	db, dbConn := setup()
	defer teardown(db, dbConn, clearChatDB)
	chatModel := &ChatModel{DB: db}
	seedChatDB(db)
	chatID := 1
	expected := []Chat{
		{
			ID:   2,
			Name: "Test Chat Room 2",
		},
	}

	err := chatModel.DeleteChat(chatID)
	result, _ := chatModel.GetAllChat()

	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(result))
}
