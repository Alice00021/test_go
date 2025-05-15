package entity

import (
	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Name     string
	Female 	 string
	Username string  
    Password string
}
