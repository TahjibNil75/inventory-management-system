package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"image/png"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/inventory-management-system/models/dto"
)

func GenerateQRCode(data dto.AssetQRCode) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("error encoding json data")
	}
	qrCode, err := qr.Encode(string(jsonData), qr.M, qr.Auto)
	if err != nil {
		return nil, errors.New("error generating qr code")
	}
	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		return nil, errors.New("error scaling QR code")
	}

	pngImgae := new(bytes.Buffer)
	err = png.Encode(pngImgae, qrCode)
	if err != nil {
		return nil, errors.New("failed to conver in PNG")
	}
	return pngImgae.Bytes(), nil
}

func UploadQRcode(data []byte, key string) error {
	sess, err := CreateS3Session()
	if err != nil {
		return errors.New("failed creating session")
	}
	s3Client := s3.New(sess)

	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("AWS_QRCODE_BUCKET")),
		Key:           aws.String(key),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
		ContentType:   aws.String("image/png"),
		// ACL:           aws.String("public-read"),
	})
	return err
}
