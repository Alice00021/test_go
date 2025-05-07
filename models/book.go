package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Book struct {
    gorm.Model
    Title      string  
    AuthorID   uint      
    Author     Author    
}

func (b *Book) BeforeCreate(tx *gorm.DB) (err error){
	var author_exist Author

	if err := tx.First(&author_exist, b.AuthorID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("автор с ID %d не найден", b.AuthorID)
        }
        return fmt.Errorf("ошибка проверки автора: %w", err)
    }
    return nil
}
