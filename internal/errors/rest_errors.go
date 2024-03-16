package rest_errors

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidEmailOrPassword = NewError("invalid email or password")
	ErrInvalidAccessToken     = NewError("invalid access token token")
	ErrCreateAccessToken      = NewError("failed to create access token")

	ErrBadRequest       = "Bad request"
	ErrInvalidParameter = "Invalid parameters"
)

func NewError(msg string) error {
	return errors.New(msg)
}

type RestErr struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
	Status  int    `json:"status"`
}

func NewBadRequestError(msg string) *RestErr {
	restErr := &RestErr{
		Message: ErrBadRequest,
		Status:  http.StatusInternalServerError,
	}
	return restErr
}
