package csv_controller

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/inventory-management-system/config"
// 	"github.com/inventory-management-system/utils"
// )

// func BackupandDownloadCSV(c *gin.Context) {

// 	db := config.ConnectToDB()
// 	csvData, err := utils.GenerateCSVFromDB(db)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "failed to generate csv from db",
// 		})
// 		return
// 	}
// 	fileName := utils.GenerateFileName()
// 	err = utils.UploadToS3(csvData, fileName)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "failed to upload in s3",
// 		})
// 		return
// 	}
// 	// Set response headers for downloading
// 	c.Header("Content-Disposition", "attachment; filename="+fileName)
// 	c.Header("Content-Type", "text/csv")

// 	c.Data(http.StatusOK, "text/csv", csvData)

// }
