package dto

import (
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserRequest struct {
	Id       int     `json:"uid" gorm:"primaryKey:unique"`
	Email    string  `gorm:"unique" json:"email"`
	Password *string `json:"password,omitempty"`
}

type UserResponse struct {
	Id        int       `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (vr UserRequest) Validate() error {
	return v.ValidateStruct(&vr,
		v.Field(&vr.Email, v.Required, is.Email),
		v.Field(&vr.Password, v.Required, v.Length(4, 10)),
	)
}
