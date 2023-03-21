package service

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	defaultPageSize   = 10
	defaultPageNumber = 1
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
	msg1 := Message{ChatID: 1, Data: "Test Message 1"}
	msg2 := Message{ChatID: 1, Data: "Test Message 2"}

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
			ID:     2,
			ChatID: 1,
			Data:   "Test Message 2",
		},
		{
			ID:     1,
			ChatID: 1,
			Data:   "Test Message 1",
		},
	}

	messages, err := messageModel.GetMessagesInChat(1, defaultPageNumber, defaultPageSize)

	assert.NoError(t, err)
	assert.Equal(t, expected[0].Data, messages[0].Data)
}

func TestPaginationQuery(t *testing.T) {
	tests := []struct {
		name       string
		pageNumber int
		pageSize   int
		expected   []Message
	}{
		{
			name:       "return only one message when pageSize is 1",
			pageNumber: 1,
			pageSize:   1,
			expected:   []Message{{ID: 2, ChatID: 1, Data: "Test Message 2"}},
		},
		{
			name:       "return two message when pageSize is 2",
			pageNumber: 1,
			pageSize:   2,
			expected:   []Message{{ID: 2, ChatID: 1, Data: "Test Message 2"}, {ID: 1, ChatID: 1, Data: "Test Message 1"}},
		},
		{
			name:       "return only item in the specified page number -1",
			pageNumber: 1,
			pageSize:   1,
			expected:   []Message{{ID: 2, ChatID: 1, Data: "Test Message 2"}},
		},
		{
			name:       "return only item in the specified page number -2",
			pageNumber: 2,
			pageSize:   1,
			expected:   []Message{{ID: 1, ChatID: 1, Data: "Test Message 1"}},
		},
		{
			name:       "return empty list when no data exist in the specified page number",
			pageNumber: 3,
			pageSize:   1,
			expected:   []Message{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, dbConn := setup()
			defer teardown(db, dbConn, clearMessageDB)
			messageModel := &MessageModel{DB: db}
			seedMessageDB(db)

			msg, err := messageModel.GetMessagesInChat(1, test.pageNumber, test.pageSize)

			assert.NoError(t, err)
			if assert.Equal(t, len(test.expected), len(msg)) && len(test.expected) != 0 {
				assert.Equal(t, test.expected[0].Data, msg[0].Data)
			}
		})
	}
}

func TestStoreNewMessage(t *testing.T) {
	db, dbConn := setup()
	defer teardown(db, dbConn, clearMessageDB)
	messageModel := &MessageModel{DB: db}
	seedMessageDB(db)
	newMsg := &Message{ChatID: 1, Data: "Test Message 3"}
	expected := []Message{
		{
			ID:     3,
			ChatID: 1,
			Data:   "Test Message 3",
		},
		{
			ID:     2,
			ChatID: 1,
			Data:   "Test Message 2",
		},
		{
			ID:     1,
			ChatID: 1,
			Data:   "Test Message 1",
		},
	}

	err := messageModel.StoreNewMessage(newMsg)

	result, _ := messageModel.GetMessagesInChat(1, defaultPageNumber, defaultPageSize)
	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(result))
	assert.Equal(t, expected[2].Data, result[2].Data)
}
