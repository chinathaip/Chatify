package service

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func seedParticipantDB(db *gorm.DB) {
	db.AutoMigrate(&ChatParticipants{})
	p1 := ChatParticipants{ChatID: 1, UserID: uuid.UUID{}}
	p2 := ChatParticipants{ChatID: 2, UserID: uuid.UUID{}}

	if err := db.Create(&p1).Error; err != nil {
		log.Fatalln("cannot seed p1: ", err)
	}

	if err := db.Create(&p2).Error; err != nil {
		log.Fatalln("cannot seed p2: ", err)
	}
}

func clearParticipantDB(db *gorm.DB) {
	if err := db.Delete(&ChatParticipants{}, "chat_participant_id>=1").Error; err != nil {
		log.Fatal(err)
	}
}
func TestGetAllParticipant(t *testing.T) {
	t.Run("Should Return all participants (2) in the table", func(t *testing.T) {
		db, dbConn := setup()
		defer teardown(db, dbConn, clearParticipantDB)
		seedParticipantDB(db)
		var (
			participantModel = &ParticipantModel{DB: db}
			expectedSize     = 2
		)

		result, err := participantModel.GetAllParticipants()

		assert.NoError(t, err)
		assert.Equal(t, expectedSize, len(result))
	})
}

func TestAddAsParticipant(t *testing.T) {
	t.Run("Should Add New Row when unique ChatID or UserID", func(t *testing.T) {
		db, dbConn := setup()
		defer teardown(db, dbConn, clearParticipantDB)
		seedParticipantDB(db)
		var (
			participantModel = &ParticipantModel{DB: db}
			newParticipant   = &ChatParticipants{ChatID: 3, UserID: uuid.UUID{}}
			expectedSize     = 3
		)

		err := participantModel.AddAsParticipant(newParticipant)
		assert.NoError(t, err)
		result, err := participantModel.GetAllParticipants()

		assert.NoError(t, err)
		assert.Equal(t, expectedSize, len(result))
		assert.Equal(t, newParticipant.ChatID, result[2].ChatID)
	})

	t.Run("Should Do Nothing when tries to input same data", func(t *testing.T) {
		db, dbConn := setup()
		defer teardown(db, dbConn, clearParticipantDB)
		seedParticipantDB(db)
		var (
			participantModel = &ParticipantModel{DB: db}
			newParticipant   = &ChatParticipants{ChatID: 1, UserID: uuid.UUID{}}
			expectedSize     = 2
		)

		err := participantModel.AddAsParticipant(newParticipant)
		assert.NoError(t, err)
		result, err := participantModel.GetAllParticipants()

		assert.NoError(t, err)
		assert.Equal(t, expectedSize, len(result))
	})
}
