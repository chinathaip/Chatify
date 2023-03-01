package service

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setup() (*gorm.DB, *sql.DB) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	dbConn, _ := db.DB()
	return db, dbConn
}

func teardown(db *gorm.DB, dbConn *sql.DB, clearTable func(db *gorm.DB)) {
	clearTable(db)
	dbConn.Close()
}

func seedMessageDB(db *gorm.DB) {
	db.AutoMigrate(&Message{})
	msg1 := Message{SenderID: "1", ChatID: 1, Data: "Test Message 1"}
	msg2 := Message{SenderID: "1", ChatID: 1, Data: "Test Message 2"}

	if err := db.Create(&msg1).Error; err != nil {
		log.Fatalln("cannot seed msg 1", err)
	}

	if err := db.Create(&msg2).Error; err != nil {
		log.Fatalln("cannot seed msg 2", err)
	}
}

func clearMessageDB(db *gorm.DB) {
	if err := db.Delete(&Message{}, "chat_id>=1").Error; err != nil {
		log.Fatal(err)
	}
}

func TestGetMessageInChat(t *testing.T) {
	db, dbConn := setup()
	defer teardown(db, dbConn, clearMessageDB)
	messageModel := &MessageModel{DB: db}
	seedMessageDB(db)
	expected := []Message{
		{
			ID:       1,
			SenderID: "1",
			ChatID:   1,
			Data:     "Test Message 1",
		},
		{
			ID:       2,
			SenderID: "1",
			ChatID:   1,
			Data:     "Test Message 2",
		},
	}

	messages, err := messageModel.GetMessagesInChat(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, messages)
}

func TestStoreNewMessage(t *testing.T) {
	db, dbConn := setup()
	defer teardown(db, dbConn, clearMessageDB)
	messageModel := &MessageModel{DB: db}
	seedMessageDB(db)
	newMsg := &Message{SenderID: "1", ChatID: 1, Data: "Test Message 3"}
	expected := []Message{
		{
			ID:       1,
			SenderID: "1",
			ChatID:   1,
			Data:     "Test Message 1",
		},
		{
			ID:       2,
			SenderID: "1",
			ChatID:   1,
			Data:     "Test Message 2",
		},
		{
			ID:       3,
			SenderID: "1",
			ChatID:   1,
			Data:     "Test Message 3",
		},
	}

	err := messageModel.StoreNewMessage(newMsg)

	result, _ := messageModel.GetMessagesInChat(1)
	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(result))
	assert.Equal(t, expected[2].Data, result[2].Data)
}
