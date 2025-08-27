package usecase

import (
	"context"
	"test_go/internal/entity"
)

type UserUseCase interface {
	Register(context.Context, *entity.UserInput) error
	Login(context.Context, string, string) (*TokenPair, error)
	RefreshTokens(context.Context, string) (*TokenPair, error)
	UpdateUser(context.Context, *entity.User) error
	GetByUserName(context.Context, string) (*entity.User, error)
	GetAllUsers(context.Context) ([]entity.User, error)
	ChangePassword(context.Context, string, string) error
	/* 	Logout(ctx context.Context, username string, password string) error */
}
