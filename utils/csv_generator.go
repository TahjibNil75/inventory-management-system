package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/inventory-management-system/models"
	"gorm.io/gorm"
)

func GenerateCSVFromDB(db *gorm.DB) ([]byte, error) {
	// Query to select data from asset_details table
	query := "SELECT * FROM asset_details"

	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, errors.New("error executing query")
	}
	defer rows.Close()

	var csvData bytes.Buffer
	writer := csv.NewWriter(&csvData)

	columnHeaders := []string{"ID", "USER NAME", "ASSET TYPE", "PRICE", "PURCHASED FROM", "PURCHASE DATE", "SERIAL NUMBER", "ASSET TAG", "MANUFACTURER", "MODEL", "OS TYPE", "LOCATION"}
	err = writer.Write(columnHeaders)
	if err != nil {
		return nil, errors.New("error writing column headers to CSV")
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
			return nil, errors.New("error scanning row")
		}
		// Write data to CSV
		err = writer.Write([]string{fmt.Sprint(id), user_name, asset_type, fmt.Sprint(price), purchased_from, purchase_date.Format("2006-01-02_15-04-05"), serial_number, asset_tag, manufacturer, model, os_type, location})
		if err != nil {
			return nil, errors.New("error writing to CSV")
		}
	}
	writer.Flush()

	return csvData.Bytes(), nil
}

func UploadToS3(data []byte, key string) error {
	sess, err := CreateS3Session()
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

func GenerateFileName() string {
	// Get the current UTC time
	now := time.Now().UTC()

	// Format the current time as "AssetList_day-month-year.csv"
	fileName := "AssetList_" + now.Format("02-January-2006") + ".csv"

	return fileName
}

func GenerateCSVFromAssetDetails(assets []models.AssetDetails) ([]byte, string, error) {
	var csvData bytes.Buffer
	writer := csv.NewWriter(&csvData)

	columnHeaders := []string{"ID", "USER NAME", "ASSET TYPE", "PRICE", "PURCHASED FROM", "PURCHASE DATE", "SERIAL NUMBER", "ASSET TAG", "MANUFACTURER", "MODEL", "OS TYPE", "LOCATION"}
	if err := writer.Write(columnHeaders); err != nil {
		return nil, "", errors.New("error writing column headers to CSV")
	}

	for _, asset := range assets {
		// Write data to CSV
		err := writer.Write([]string{
			fmt.Sprint(asset.Id),
			asset.UserName,
			asset.AssetType,
			fmt.Sprint(asset.Price),
			asset.PurchasedFrom,
			asset.PurchaseDate.Format("2006-01-02_15-04-05"),
			asset.SerialNumber,
			asset.AssetTag,
			asset.Manufacturer,
			asset.Model,
			asset.OsType,
			asset.Location,
		})
		if err != nil {
			return nil, "", errors.New("error writing to CSV")
		}
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, "", err
	}

	fileName := GenerateFileName()

	return csvData.Bytes(), fileName, nil
}
