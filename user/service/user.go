package service

import (
	"github.com/inventory-management-system/models"
	"github.com/inventory-management-system/models/dto"
	"github.com/inventory-management-system/user/repository"
	"github.com/inventory-management-system/utils"
)

type User struct {
	repo repository.UserRepository
}

type UserService interface {
	CreateUser(user dto.UserRequest) (*models.User, error)
}

func NewUserService(customRepo repository.UserRepository) UserService {
	return &User{
		repo: customRepo,
	}
}

func (usr *User) CreateUser(user dto.UserRequest) (*models.User, error) {
	usrModel := models.User{}
	_ = utils.StructToStruct(user, &usrModel)
	hashPass, _ := utils.GenerateHash(usrModel.Password)
	usrModel.Password = hashPass
	resp, err := usr.repo.CreateUserRepository(usrModel)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
