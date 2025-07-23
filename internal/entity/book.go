package entity

import (
	/* "errors"
	"fmt" */

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title    string
	AuthorID uint
	Author   Author
	pageSize uint
}
