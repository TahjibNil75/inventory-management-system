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
	grp.PUT("v1/asset-details/update/:asset_id", ctl.UpdateAsset)
	grp.GET("v1/list-of-assets", ctl.ListOfAssets)
	grp.DELETE("v1/asset-details/delete/:asset_id", ctl.DeleteAsset)
	grp.GET("v1/asset/:asset_id", ctl.GetAssetById)
	grp.POST("v1/asset-details/csv", ctl.ExportCSV)
	grp.POST("v1/asset-detail/qrcode/:asset_id", ctl.AssetQRCode)
	grp.GET("v1/asset-details/search", ctl.SearchByKeyWord)

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
	assetId := c.Param("asset_id")
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
	assetId := c.Param("asset_id")
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

func (ctl *Asset) GetAssetById(c *gin.Context) {
	assetID := c.Param("asset_id")

	id, err := strconv.Atoi(assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to convert id to int",
		})
		return
	}
	resp, err := ctl.svc.GetAssetsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get asset",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": resp,
	})

}

func (ctl *Asset) ExportCSV(c *gin.Context) {
	assets, err := ctl.svc.GetAllAssets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve asset details",
		})
		return
	}
	csvData, fileName, err := utils.GenerateCSVFromAssetDetails(assets)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate CSV",
		})
		return
	}
	// Upload CSV to S3
	err = utils.UploadToS3(csvData, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to upload to S3",
		})
		return
	}
	// Set response headers for downloading
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "text/csv")
	// Download CSV from S3 and return to client
	c.Data(http.StatusOK, "text/csv", csvData)
}

func (ctl *Asset) AssetQRCode(c *gin.Context) {
	assetID := c.Param("asset_id")
	id, err := strconv.Atoi(assetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to convert id to int",
		})
		return
	}
	resp, err := ctl.svc.GetAssetsByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "failed to get asset or internal server error",
		})
		return
	}
	qrCodeData := dto.AssetQRCode{
		UserName:      resp.UserName,
		AssetTag:      resp.AssetTag,
		SerialNumber:  resp.SerialNumber,
		AssetType:     resp.AssetType,
		Location:      resp.Location,
		PurchasedFrom: resp.PurchasedFrom,
	}
	qrCodeBytes, err := utils.GenerateQRCode(qrCodeData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to generate QR code",
		})
		return
	}
	err = utils.UploadQRcode(qrCodeBytes, "qrcode/"+qrCodeData.AssetTag+".png")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to upload qrcode in s3",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "qrcode was successfully uploaded to s3",
	})

	c.JSON(http.StatusOK, gin.H{
		"qrCode": qrCodeBytes,
	})
}

func (ctl *Asset) SearchByKeyWord(c *gin.Context) {
	keyWord := c.Query("keyword")
	if keyWord == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "keyword is required",
		})
		return
	}
	assets, err := ctl.svc.SearchAssetByKeyWord(keyWord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to search assets by keyword",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": assets,
	})
}
