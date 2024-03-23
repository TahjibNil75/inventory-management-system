package csvgenerator

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"github.com/inventory-management-system/config"
	"github.com/inventory-management-system/utils"
	"gorm.io/gorm"

	"github.com/aws/aws-sdk-go/service/s3"
)

func uploadToS3(data []byte, key string) error {
	sess, err := utils.CreateS3Session()
	if err != nil {
		return errors.New("failed creating session")
	}
	s3Client := s3.New(sess)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("AWS_BUCKET")),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
		ContentType:   aws.String("text/csv"),
		// ACL:           aws.String("public-read"),
	})
	return err
}

func generateAndUploadCSV(db *gorm.DB) error {

	// Query to select data from asset_details table
	query := "SELECT * FROM asset_details"

	// Execute the query
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var csvData bytes.Buffer
	writer := csv.NewWriter(&csvData)

	columnHeaders := []string{"ID", "USER NAME", "ASSET TYPE", "PRICE", "PURCHASED FROM", "PURCHASE DATE", "SERIAL NUMBER", "ASSET TAG", "MANUFACTURER", "MODEL", "OS TYPE", "LOCATION"}
	err = writer.Write(columnHeaders)
	if err != nil {
		return fmt.Errorf("error writing column headers to CSV: %w", err)
	}
	for rows.Next() {
		var (
			id             int
			user_name      string
			asset_type     string
			price          float64
			status         bool
			purchased_from string
			purchase_date  time.Time
			serial_number  string
			asset_tag      string
			manufacturer   string
			model          string
			os_type        string
			location       string
		)
		err := rows.Scan(&id, &user_name, &asset_type, &price, &status, &purchased_from, &purchase_date, &serial_number, &asset_tag, &manufacturer, &model, &os_type, &location)
		if err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}
		// Write data to CSV
		err = writer.Write([]string{fmt.Sprint(id), user_name, asset_type, fmt.Sprint(price), purchased_from, purchase_date.Format("2006-01-02_15-04-05"), serial_number, asset_tag, manufacturer, model, os_type, location})
		if err != nil {
			return fmt.Errorf("error writing to CSV: %w", err)
		}
	}
	writer.Flush()

	now := time.Now().UTC().Format("2006-01-02_15-04-05")

	err = uploadToS3(csvData.Bytes(), fmt.Sprintf("asset_details_%s.csv", now))
	if err != nil {
		return fmt.Errorf("error uploading to S3: %w", err)
	}

	return nil
}

func UploadCSV(c *gin.Context) {
	db := config.ConnectToDB()

	err := generateAndUploadCSV(db)
	if err != nil {
		fmt.Println("Error generating and uploading CSV:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to generate and upload CSV to S3",
		})
		return
	}
	fmt.Println("CSV uploaded to S3 successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "CSV uploaded to S3 successfully",
	})

}
