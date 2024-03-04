package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/inventory-management-system/models/dto"
	"github.com/inventory-management-system/user/service"
	"github.com/inventory-management-system/utils"
)

type User struct {
	UserService service.UserService
}

func NewUserController(g interface{}, cusSvc service.UserService) {
	usr := &User{
		UserService: cusSvc,
	}
	group := g.(*gin.RouterGroup)
	group.POST("/v1/user/signup", usr.SignUp)
}

func (ctr *User) SignUp(ctx *gin.Context) {
	reqBody := dto.UserRequest{}
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "binding json failed",
		})
		return
	}
	// Validate Request
	if err := reqBody.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := ctr.UserService.CreateUser(reqBody)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
		ctx.JSON(http.StatusInternalServerError, err)
	}
	userResp := dto.UserResponse{
		Id:        resp.Id,
		Email:     reqBody.Email,
		CreatedAt: resp.CreatedAt,
	}
	_ = utils.StructToStruct(*resp, userResp)
	ctx.JSON(http.StatusOK, gin.H{
		"user": userResp,
	})
}
