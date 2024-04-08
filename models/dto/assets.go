package dto

import (
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
)

type AssetEntryReq struct {
	UserName      string    `json:"user_name"`
	AssetType     string    `json:"asset_type"`
	Price         float64   `json:"price"`
	Status        string    `json:"status"`
	PurchasedFrom string    `json:"purchased_from"`
	PurchaseDate  time.Time `json:"purchase_date"`
	SerialNumber  string    `json:"serial_number"`
	AssetTag      string    `json:"asset_tag"`
	Manufacturer  string    `json:"manufacturer"`
	Model         string    `json:"model"`
	OsType        string    `json:"os_type"`
	Location      string    `json:"location"`
}

type AssetUpdateReq struct {
	UserName      string    `json:"user_name"`
	AssetType     string    `json:"asset_type"`
	Price         float64   `json:"price"`
	Status        string    `json:"status"`
	PurchasedFrom string    `json:"purchased_from"`
	PurchaseDate  time.Time `json:"purchase_date"`
	SerialNumber  string    `json:"serial_number"`
	AssetTag      string    `json:"asset_tag"`
	Manufacturer  string    `json:"manufacturer"`
	Model         string    `json:"model"`
	OsType        string    `json:"os_type"`
	Location      string    `json:"location"`
}

type AssetDetailsResp struct {
	Id            int       `json:"id"`
	UserName      string    `json:"user_name"`
	AssetType     string    `json:"asset_type"`
	Price         float64   `json:"price"`
	Status        string    `json:"status"`
	PurchasedFrom string    `json:"purchased_from"`
	PurchaseDate  time.Time `json:"purchase_date"`
	SerialNumber  string    `json:"serial_number"`
	AssetTag      string    `json:"asset_tag"`
	Manufacturer  string    `json:"manufacturer"`
	Model         string    `json:"model"`
	OsType        string    `json:"os_type"`
	Location      string    `json:"location"`
	CreatedAt     time.Time `json:"created_at"`
}

type AssetQRCode struct {
	UserName      string
	AssetType     string
	AssetTag      string
	SerialNumber  string
	Price         float64
	PurchasedFrom string
	PurchaseDate  time.Time
	Location      string
}

func (a AssetEntryReq) Validate() error {
	return v.ValidateStruct(&a,
		v.Field(&a.UserName, v.Required),
		v.Field(&a.AssetType, v.Required),
		v.Field(&a.Price, v.Required),
		v.Field(&a.Status, v.Required),
		v.Field(&a.PurchasedFrom, v.Required),
		v.Field(&a.PurchaseDate, v.Required),
		v.Field(&a.SerialNumber, v.Required),
		v.Field(&a.AssetTag, v.Required),
		v.Field(&a.Manufacturer, v.Required),
		v.Field(&a.Model, v.Required),
		v.Field(&a.OsType, v.Required),
		v.Field(&a.Location, v.Required),
	)
}

func (a AssetUpdateReq) Validate() error {
	return v.ValidateStruct(&a,
		v.Field(&a.UserName, v.Required),
		v.Field(&a.AssetType, v.Required),
		v.Field(&a.Price, v.Required),
		v.Field(&a.Status, v.Required),
		v.Field(&a.PurchasedFrom, v.Required),
		v.Field(&a.PurchaseDate, v.Required),
		v.Field(&a.SerialNumber, v.Required),
		v.Field(&a.AssetTag, v.Required),
		v.Field(&a.Manufacturer, v.Required),
		v.Field(&a.Model, v.Required),
		v.Field(&a.OsType, v.Required),
		v.Field(&a.Location, v.Required),
	)
}
