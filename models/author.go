package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Author struct {
    gorm.Model
    Name     string
	Female 	 bool  
    Books    []Book    
}

func (a *Author) BeforeCreate(tx *gorm.DB) (err error) {
	 var existing_author Author 
  
	 if err := tx.Where("name = ?", a.Name).First(&existing_author).Error; err != nil {
		return fmt.Errorf("автор с таким именем уже существует")
	}
	
	return nil
}
