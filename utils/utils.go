package utils

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

func StructToStruct(input interface{}, output interface{}) error {
	if b, err := json.Marshal(input); err == nil {
		return json.Unmarshal(b, &output)
	} else {
		return err
	}
}

func GenerateHash(password *string) (*string, error) {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(*password), 8)
	hashPassword := string(hashedPass)
	return &hashPassword, nil
}
