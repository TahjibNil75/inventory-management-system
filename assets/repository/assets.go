package asset_repository

import (
	"errors"

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
	DeleteAssetById(assetId int) error
	GetAssetById(assetId int) (*models.AssetDetails, error)
	SearchByKeyWord(keyword string) ([]models.AssetDetails, error)
}

func NewAssetRepository(db *gorm.DB) AssetRepo {
	return &Asset{
		DB: db,
	}
}

func (assetRepo *Asset) CreateAsset(asset models.AssetDetails) (*models.AssetDetails, error) {
	err := assetRepo.DB.Model(&models.AssetDetails{}).Create(&asset).Error
	if err != nil {
		return nil, errors.New("error creating")
	}
	return &asset, nil
}

func (assetRepo *Asset) GetAssetById(assetId int) (*models.AssetDetails, error) {
	resp := models.AssetDetails{}
	err := assetRepo.DB.Model(&models.AssetDetails{}).Where("id = ?", assetId).Find(&resp).Error
	if err != nil {
		return nil, errors.New("error ")
	}
	if resp.Id == 0 {
		return nil, errors.New("invalid asset id")
	}
	return &resp, nil
}

func (assetRepo *Asset) UpdateAssetById(asset models.AssetDetails) (*models.AssetDetails, error) {
	err := assetRepo.DB.Model(&models.AssetDetails{}).Where("id = ?", asset.Id).Updates(&asset).Error
	if err != nil {
		return nil, errors.New("error updating")
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

func (assetRepo *Asset) DeleteAssetById(assetId int) error {
	err := assetRepo.DB.Model(&models.AssetDetails{}).Where("id = ?", assetId).Delete(&models.AssetDetails{}).Error
	if err != nil {
		return nil
	}
	return nil
}

func (assetRepo *Asset) SearchByKeyWord(keyword string) ([]models.AssetDetails, error) {
	var allAssets []models.AssetDetails
	err := assetRepo.DB.Where("user_name LIKE ? OR asset_type LIKE ? OR purchased_from LIKE ? OR serial_number LIKE ? OR asset_tag LIKE ? OR manufacturer LIKE ? OR model LIKE ? OR os_type LIKE ? OR location LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Find(&allAssets).Error
	if err != nil {
		return nil, errors.New("error searching")
	}
	return allAssets, nil
}
