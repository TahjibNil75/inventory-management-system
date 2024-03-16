package asset_repository

import (
	"github.com/inventory-management-system/models"
	"gorm.io/gorm"
)

type Asset struct {
	*gorm.DB
}

type AssetRepo interface {
	CreateAsset(asset models.AssetDetails) (*models.AssetDetails, error)
	UpdateAssetById(asset models.AssetDetails) (*models.AssetDetails, error)
	GetAllAssets() ([]models.AssetDetails, error)
}

func NewAssetRepository(db *gorm.DB) AssetRepo {
	return &Asset{
		DB: db,
	}
}

func (assetRepo *Asset) CreateAsset(asset models.AssetDetails) (*models.AssetDetails, error) {
	err := assetRepo.DB.Model(&models.AssetDetails{}).Create(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (assetRepo *Asset) UpdateAssetById(asset models.AssetDetails) (*models.AssetDetails, error) {
	err := assetRepo.DB.Model(&models.AssetDetails{}).Where("id = ?", asset.Id).Updates(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (assetRepo *Asset) GetAllAssets() ([]models.AssetDetails, error) {
	var allAssets []models.AssetDetails
	err := assetRepo.DB.Model(&models.AssetDetails{}).Find(&allAssets).Error
	if err != nil {
		return nil, err
	}
	return allAssets, nil
}
