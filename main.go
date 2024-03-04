package main

import (
	"github.com/gin-gonic/gin"
	"github.com/inventory-management-system/config"
	"github.com/inventory-management-system/user/controller"
	"github.com/inventory-management-system/user/repository"
	"github.com/inventory-management-system/user/service"
)

func main() {
	// Connect to the database
	db := config.ConnectToDB()

	// Run database migrations
	config.Migrate(db)

	// Initialize Gin router
	router := gin.Default()

	// Initialize user repository
	userRepo := repository.NewUserRepository(db)

	// Initialize user service
	userService := service.NewUserService(userRepo)

	// Register user controller
	userGroup := router.Group("")
	controller.NewUserController(userGroup, userService)

	// Start the server
	router.Run(":8080")
}
