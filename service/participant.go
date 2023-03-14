package service

import (
	"gorm.io/gorm"
)

type ParticipantService interface {
	AddAsParticipant(*ChatParticipants) error
}

type ParticipantModel struct {
	DB *gorm.DB
}

func (m *ParticipantModel) AddAsParticipant(p *ChatParticipants) error {

	//check if already exist
	existing := ChatParticipants{}
	if err := m.DB.Where(&ChatParticipants{ChatID: p.ChatID, UserID: p.UserID}).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	//if already exist, do nothing
	if existing.ID != 0 || existing.ChatID != 0 {
		return nil
	}

	//insert new one if not
	if err := m.DB.Create(p).Error; err != nil {
		return err
	}
	return nil
}
