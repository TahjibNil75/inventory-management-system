package main

import (
	"github.com/gin-gonic/gin"
	asset_controller "github.com/inventory-management-system/assets/controller"
	asset_repository "github.com/inventory-management-system/assets/repository"
	asset_service "github.com/inventory-management-system/assets/service"
	auth "github.com/inventory-management-system/auth/controller"
	"github.com/inventory-management-system/auth/middleware"
	authservice "github.com/inventory-management-system/auth/service"
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
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	// Initialize user repository
	userRepo := repository.NewUserRepository(db)
	// Initialize user service
	userService := service.NewUserService(userRepo)
	// Register user controller
	userGroup := router.Group("")
	controller.NewUserController(userGroup, userService)

	// Initialize auth service
	authSvc := authservice.NewAuthService(userRepo)

	// Register auth controller
	authGroup := router.Group("")
	auth.NewAuthController(authGroup, authSvc)

	// Use middleware to ensure authentication
	router.Use(middleware.AuthMiddleware())

	// Initialize asset repository
	assetRepo := asset_repository.NewAssetRepository(db)
	// Initialize asset service
	assetService := asset_service.NewAssetService(assetRepo)
	// Register asset controller
	assetGroup := router.Group("api")
	asset_controller.NewAssetController(assetGroup, assetService)

	// router.POST("/api/assets_details/csv", csv_controller.BackupandDownloadCSV)

	// Start the server
	router.Run(":8080")
}
