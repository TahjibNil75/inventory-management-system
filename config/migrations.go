package config

import (
	"github.com/inventory-management-system/models"
	"gorm.io/gorm"
)

func Migrate(DB *gorm.DB) {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.AssetDetails{})
}
