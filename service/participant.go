package service

import (
	"gorm.io/gorm"
)

type ParticipantService interface {
	AddAsParticipant(*ChatParticipants) error
	Exist(*ChatParticipants) bool
}

type ParticipantModel struct {
	DB *gorm.DB
}

func (m *ParticipantModel) AddAsParticipant(p *ChatParticipants) error {
	if yes := m.Exist(p); yes {
		return nil
	}

	//insert new one if not
	if err := m.DB.Create(p).Error; err != nil {
		return err
	}

	return nil
}

func (m *ParticipantModel) GetAllParticipants() ([]ChatParticipants, error) {
	var cp []ChatParticipants
	if err := m.DB.Find(&cp).Error; err != nil {
		return nil, err
	}

	return cp, nil
}

func (m *ParticipantModel) Exist(p *ChatParticipants) bool {
	existing := ChatParticipants{}
	if err := m.DB.Where(&ChatParticipants{ChatID: p.ChatID, UserID: p.UserID}).First(&existing).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return true
		}
	}

	//if already exist, do nothing
	if existing.ID != 0 || existing.ChatID != 0 {
		return true
	}

	return false
}
