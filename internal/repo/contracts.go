/* интерфейсы для репозиториев */
package repo

import (
	"context"
	"test_go/internal/entity"
)

type BookRepository interface {
	Create(context.Context, *entity.Book) error
	Update(context.Context, *entity.Book, uint) error
	Delete(context.Context, uint) error
	GetByID(context.Context, uint) (*entity.Book, error)
	GetAll(context.Context) ([]entity.Book, error)
}

type AuthorRepository interface {
	Create(context.Context, *entity.Author) error
	Update(context.Context, *entity.Author, uint) error
	Delete(context.Context, uint) error
	GetByID(context.Context, uint) (*entity.Author, error)
	GetAll(context.Context) ([]entity.Author, error)
}

type UserRepository interface {
	Create(context.Context, *entity.User) error
	Update(context.Context, *entity.User) error
	GetByUserName(context.Context, string) (*entity.User, error)
	GetById(context.Context, uint) (*entity.User, error)
	GetAll(context.Context) ([]entity.User, error)
	GetByEmail(context.Context, string) (*entity.User, error)
	GetByVerifyToken(context.Context, string) (*entity.User, error)
	UpdateRating(context.Context, uint, float32) error
}
