package models

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// User struct
type User struct {
	CommonModelFields
	Name     string `gorm:"not null" json:"name,omitempty"`
	Email    string `gorm:"unique;not null" json:"email,omitempty"`
	Username string `gorm:"unique;not null" json:"username,omitempty"`
	Password string `gorm:"not null" json:"password,omitempty"`
	Phone    string `gorm:"unique" json:"phone,omitempty" sql:"DEFAULT:NULL"`
}

// BeforeSave -> hook for hashing password before save to DB
func (u *User) BeforeSave() error {
	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// MarshalJSON -> remove password from get user data
func (u User) MarshalJSON() ([]byte, error) {
	type user User // prevent recursion
	x := user(u)
	x.Password = ""
	return json.Marshal(x)
}
