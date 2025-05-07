package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Book struct {
    gorm.Model
    Title      string  
    AuthorID   uint      
    Author     Author    
}

func (b *Book) BeforeCreate(db *gorm.DB) (err error){
	var author_exist Author

	if err := db.First(&author_exist, b.AuthorID).Error; err != nil {
        return fmt.Errorf("автор не найден")
    }
	return nil
}
