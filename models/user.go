package models

import "time"

type User struct {
	Id        int        `json:"id"`
	Email     string     `gorm:"unique" json:"email"`
	Password  *string    `json:"password,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
