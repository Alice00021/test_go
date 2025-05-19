/* интерфейсы для репозиториев */
package repo

import (
	"context"
	"test_go/internal/entity"
)

type BookRepository interface{
	Create(ctx context.Context, book *entity.Book) error
    Update(ctx context.Context, book *entity.Book, id uint) error
    Delete(ctx context.Context, id uint) error
    GetByID(ctx context.Context, id uint) (*entity.Book, error)
    GetAll(ctx context.Context) ([]entity.Book, error)
}

type AuthorRepository interface {
    Create(ctx context.Context, author *entity.Author) error
    Update(ctx context.Context, author *entity.Author, id uint) error
    Delete(ctx context.Context, id uint) error
    GetByID(ctx context.Context, id uint) (*entity.Author, error)
    GetAll(ctx context.Context) ([]entity.Author, error)
}

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    Update(ctx context.Context, user *entity.User, id uint) error
    GetByUserName(ctx context.Context, username string) (*entity.User, error)
    GetAll(ctx context.Context) ([]entity.User, error)
}
