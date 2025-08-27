package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Surname  string
	Username string
	Password string
	FilePath *string
}

type UserInput struct {
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func NewUser(name string, surname string, username string, password string) *User {
	return &User{
		Name:     name,
		Surname:  surname,
		Username: username,
		Password: password,
	}
}
