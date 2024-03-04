package repository

import (
	"errors"

	"github.com/inventory-management-system/models"
	"gorm.io/gorm"
)

type User struct {
	*gorm.DB
}

type UserRepository interface {
	CreateUserRepository(user models.User) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &User{
		DB: db,
	}
}

func (userRep *User) CreateUserRepository(user models.User) (*models.User, error) {
	err := userRep.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &user, nil
}
