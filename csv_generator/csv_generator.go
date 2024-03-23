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
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/inventory-management-system/config"
	"gorm.io/gorm"

	"github.com/aws/aws-sdk-go/service/s3"
)

func uploadToS3(data []byte, key string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		return errors.New("failed to create session")
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
	db := config.ConnectToDB() // Establish database connection

	err := generateAndUploadCSV(db) // Pass the database connection to generateAndUploadCSV
	if err != nil {
		fmt.Println("Error generating and uploading CSV:", err) // Log the error
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to generate and upload CSV to S3",
		})
		return
	}
	fmt.Println("CSV uploaded to S3 successfully") // Log success
	c.JSON(http.StatusOK, gin.H{
		"message": "CSV uploaded to S3 successfully",
	})
}
