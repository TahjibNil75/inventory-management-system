package models

import "time"

type AssetDetails struct {
	Id            int       `json:"id" gorm:"unique"`
	UserName      string    `json:"user_name"`
	AssetType     string    `json:"asset_type"`
	Price         float64   `json:"price"`
	Status        bool      `json:"status"`
	PurchasedFrom string    `json:"purchased_from"`
	PurchaseDate  time.Time `json:"purchase_date"`
	SerialNumber  string    `json:"serial_number"`
	AssetTag      string    `json:"asset_tag"`
	Manufacturer  string    `json:"manufacturer"`
	Model         string    `json:"model"`
	OsType        string    `json:"os_type"`
	Location      string    `json:"location"`
}
