package models

import (
	"gorm.io/gorm"
)

type Book struct {
    gorm.Model
    Title      string  
    AuthorID   uint      
    Author     Author    
}

type Author struct {
    gorm.Model
    Name     string
	Female 	 bool  
    Books    []Book    
}