package service

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	GetUserNameByID(uuid.UUID) (string, error)
	CreateNewUser(user User) error
}

type UserModel struct {
	DB *gorm.DB
}

func (m *UserModel) GetUserNameByID(id uuid.UUID) (string, error) {
	user := User{}
	if err := m.DB.Table("users").Select("username").Where("id=?", id).First(&user).Error; err != nil {
		return "", err
	}

	return user.Username, nil
}

func (m *UserModel) CreateNewUser(user User) error {
	if err := m.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}
