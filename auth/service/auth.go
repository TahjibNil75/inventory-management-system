package authservice

import (
	rest_errors "github.com/inventory-management-system/internal/errors"
	"github.com/inventory-management-system/models"
	"github.com/inventory-management-system/models/dto"
	"github.com/inventory-management-system/user/repository"
	"github.com/inventory-management-system/utils"
	"golang.org/x/crypto/bcrypt"
)

type IAuth interface {
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type auth struct {
	urepo repository.UserRepository
}

func NewAuthService(uRepo repository.UserRepository) IAuth {
	return &auth{
		urepo: uRepo,
	}
}

func (as *auth) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var user *models.User
	var err error

	// Retrieve user from repository
	if user, err = as.urepo.GetUserByEmail(req.Email); err != nil {
		return nil, rest_errors.ErrInvalidEmailOrPassword
	}

	// Compare the provided password with the hashed password stored in the database
	loginPass := []byte(req.Password)
	hashedPass := []byte(*user.Password)
	if err := bcrypt.CompareHashAndPassword(hashedPass, loginPass); err != nil {
		return nil, rest_errors.ErrInvalidEmailOrPassword
	}

	// Generate JWT token
	token, err := utils.CreateToken(*user)
	if err != nil {
		return nil, err
	}

	// Return the generated token in the response
	return &dto.LoginResponse{
		AccessToken: token,
	}, nil
}
