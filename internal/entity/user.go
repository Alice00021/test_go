package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string
	Surname     string
	Username    string
	Password    string
	Email       string
	IsVerified  bool
	VerifyToken *string
	FilePath    *string
}

type UserInput struct {
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

func NewUser(name, surname, username, password, email string) *User {
	return &User{
		Name:       name,
		Surname:    surname,
		Username:   username,
		Password:   password,
		Email:      email,
		IsVerified: false,
	}
}

type EmailConfig struct {
	SMTPHost       string
	SMTPPort       int
	SenderEmail    string
	SenderPassword string
}
