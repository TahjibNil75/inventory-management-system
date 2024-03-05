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
	GetUserByEmail(email string) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &User{
		DB: db,
	}
}

func (userRepo *User) CreateUserRepository(user models.User) (*models.User, error) {
	err := userRepo.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &user, nil
}

func (userRepo *User) GetUserByID(id string) (*models.User, error) {
	resp := models.User{}
	err := userRepo.DB.Model(&models.User{}).Where("id = ?", id).Find(&resp).Error
	if err != nil {
		return nil, errors.New(err.Error())
	}
	// if resp.Id == 0 {
	// 	return nil, errors.New(err.Error())
	// }
	return &resp, nil
}

func (userRepo *User) GetUserByEmail(email string) (*models.User, error) {
	resp := models.User{}
	err := userRepo.DB.Model(&models.User{}).Where("email = ?", email).Find(&resp).Error
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return &resp, nil
}
