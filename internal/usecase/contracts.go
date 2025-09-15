package usecase

import (
	"context"
	"test_go/internal/entity"
)

type (
	UserUseCase interface {
		Register(context.Context, *entity.UserInput) error
		Login(context.Context, string, string) (*TokenPair, error)
		RefreshTokens(context.Context, string) (*TokenPair, error)
		UpdateUser(context.Context, *entity.User) error
		GetByUserName(context.Context, string) (*entity.User, error)
		GetAllUsers(context.Context) ([]entity.User, error)
		ChangePassword(context.Context, string, string, string, string) error
		VerifyEmail(context.Context, string) error
		/* 	Logout(ctx context.Context, username string, password string) error */
	}
	BookUseCase interface {
		CreateBook(ctx context.Context, title string, authorID uint) (*entity.Book, error)
		UpdateBook(ctx context.Context, title string, authorID uint, id uint) (*entity.Book, error)
		DeleteBook(ctx context.Context, id uint) error
		GetByIDBook(ctx context.Context, id uint) (*entity.Book, error)
		GetAllBooks(ctx context.Context) ([]entity.Book, error)
	}
)
