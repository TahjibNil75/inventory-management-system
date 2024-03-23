package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	authservice "github.com/inventory-management-system/auth/service"
	rest_errors "github.com/inventory-management-system/internal"
	"github.com/inventory-management-system/models/dto"
)

type auth struct {
	authSvc authservice.IAuth
}

func NewAuthController(g interface{}, authSvc authservice.IAuth) {
	ath := &auth{
		authSvc: authSvc,
	}
	grp := g.(*gin.RouterGroup)
	grp.POST("/user/login", ath.Login)
}

func (ctr *auth) Login(ctx *gin.Context) {
	var cred *dto.LoginRequest
	var resp *dto.LoginResponse
	var err error

	if err := ctx.ShouldBind(&cred); err != nil {
		bodyErr := rest_errors.NewBadRequestError("Failed to parse request body")
		ctx.JSON(bodyErr.Status, bodyErr)
		return
	}
	cred.Email = strings.TrimSpace(cred.Email)
	cred.Password = strings.TrimSpace(cred.Password)

	if payloadErr := cred.Validate(); payloadErr != nil {
		restErr := rest_errors.NewBadRequestError("Failed to validate")
		ctx.JSON(restErr.Status, restErr)
		return
	}

	if resp, err = ctr.authSvc.Login(cred); err != nil {
		switch err {
		case rest_errors.ErrInvalidEmailOrPassword:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid email or password",
			})
			return
		case rest_errors.ErrCreateAccessToken:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create access token",
			})
			return
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"auth response": resp,
	})
}
