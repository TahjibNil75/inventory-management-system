package asset_service

import (
	"errors"

	asset_repository "github.com/inventory-management-system/assets/repository"
	"github.com/inventory-management-system/models"
	"github.com/inventory-management-system/models/dto"
	"github.com/inventory-management-system/utils"
)

type Asset struct {
	repo asset_repository.AssetRepo
}

type AssetService interface {
	CreateAsset(asset dto.AssetEntryReq) (*models.AssetDetails, error)
	UpdateAssetById(assetId int, asset dto.AssetUpdateReq) (*models.AssetDetails, error)
	GetAllAssets() ([]models.AssetDetails, error)
	DeleteAssetById(assetId int) error
}

func NewAssetService(assetRepo asset_repository.AssetRepo) AssetService {
	return &Asset{
		repo: assetRepo,
	}
}

func (assetSvc *Asset) CreateAsset(asset dto.AssetEntryReq) (*models.AssetDetails, error) {
	newAsset := models.AssetDetails{}
	_ = utils.StructToStruct(asset, &newAsset)
	resp, err := assetSvc.repo.CreateAsset(newAsset)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (assetSvc *Asset) UpdateAssetById(assetId int, asset dto.AssetUpdateReq) (*models.AssetDetails, error) {
	updatedAsset := models.AssetDetails{}
	_ = utils.StructToStruct(asset, &updatedAsset)
	updatedAsset.Id = assetId
	resp, err := assetSvc.repo.UpdateAssetById(updatedAsset)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (assetSvc *Asset) GetAllAssets() ([]models.AssetDetails, error) {
	resp, err := assetSvc.repo.GetAllAssets()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (assetSvc *Asset) DeleteAssetById(assetId int) error {
	_, err := assetSvc.repo.GetAssetById(assetId)
	if err != nil {
		return errors.New("asset not found or invalid asset id")
	}
	if err := assetSvc.repo.DeleteAssetById(assetId); err != nil {
		return errors.New("failed to delete asset")
	}
	return nil
}
