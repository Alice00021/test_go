package entity

import (
	/* "errors"
	"fmt" */

	"gorm.io/gorm"
)

type Author struct {
    gorm.Model
    Name     string
	Gender 	 bool  
    Books    []Book    
}

/* func (a *Author) BeforeCreate(tx *gorm.DB) (err error) {
	 var existing_author Author 
  
	err = tx.Where("name = ?", a.Name).First(&existing_author).Error

	 if err == nil {
        return fmt.Errorf("автор с именем '%s' уже существует", a.Name)
    }
	if !errors.Is(err, gorm.ErrRecordNotFound) {
        return fmt.Errorf("ошибка проверки уникальности автора: %w", err)
    }
    
    return nil
}

func (a *Author) BeforeUpdate(tx *gorm.DB) (err error){
	var original Author
	if err:= tx.First(&original, a.ID).Error; err !=nil{
		return err
	}

	if a.Name == original.Name && a.Gender == original.Gender{
		return errors.New("нет изменений")
	}
	return nil
} */