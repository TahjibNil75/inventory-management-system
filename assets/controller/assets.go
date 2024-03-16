package asset_controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	asset_service "github.com/inventory-management-system/assets/service"
	"github.com/inventory-management-system/models/dto"
	"github.com/inventory-management-system/utils"
)

type Asset struct {
	svc asset_service.AssetService
}

func NewAssetController(g interface{}, assetSvc asset_service.AssetService) {
	ctl := &Asset{
		svc: assetSvc,
	}
	grp := g.(*gin.RouterGroup)
	grp.POST("v1/asset-details/create", ctl.CreateAsset)
	grp.PUT("v1/asset-details/update/:id", ctl.UpdateAsset)
	grp.GET("v1/list-of-assets", ctl.ListOfAssets)
	grp.DELETE("v1/asset-details/delete/:id", ctl.DeleteAsset)

}

func (ctl *Asset) CreateAsset(c *gin.Context) {
	reqBody := dto.AssetEntryReq{}
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "failed to bind JSON request",
		})
		return
	}
	// Validate the request
	if payloadErr := reqBody.Validate(); payloadErr != nil {
		log.Fatal("failed to validate request body", payloadErr)
	}

	resp, err := ctl.svc.CreateAsset(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create asset",
		})
		return
	}
	assetResp := dto.AssetDetailsResp{}
	_ = utils.StructToStruct(*resp, &assetResp)
	c.JSON(http.StatusOK, gin.H{
		"AssetDetails": assetResp,
	})
}
func (ctl *Asset) UpdateAsset(c *gin.Context) {
	reqBody := dto.AssetUpdateReq{}
	if err := c.BindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "failed to bind JSON request",
		})
		return
	}
	if payloadErr := reqBody.Validate(); payloadErr != nil {
		log.Fatal("failed to validate request body", payloadErr)
	}
	assetId := c.Param("id")
	id, err := strconv.Atoi(assetId)
	if err != nil {
		log.Println("failed to convert asset id to int")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal  error",
		})
		return
	}
	resp, err := ctl.svc.UpdateAssetById(id, reqBody)
	if err != nil {
		log.Fatalf("failed to update asset")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Asset Updated": resp,
	})

}

func (ctl *Asset) ListOfAssets(c *gin.Context) {
	assets, err := ctl.svc.GetAllAssets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"List of Assets": assets,
	})
}

func (ctl *Asset) DeleteAsset(c *gin.Context) {
	assetId := c.Param("id")
	id, err := strconv.Atoi(assetId)
	if err != nil {
		log.Panicln("failed to convert asset id to int")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal error",
		})
		return
	}
	if err := ctl.svc.DeleteAssetById(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "asset deleted successfully",
	})

}
