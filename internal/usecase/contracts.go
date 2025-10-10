package usecase

import (
	"context"
	"github.com/xuri/excelize/v2"
	"mime/multipart"
	"test_go/internal/entity"
)

type (
	Auth interface {
		Register(context.Context, entity.CreateUserInput) (*entity.User, error)
		Login(context.Context, string, string) (*entity.TokenPair, error)
		VerifyEmail(context.Context, string) error
		RefreshTokens(context.Context, string) (*entity.TokenPair, error)
		Validation(context.Context, string) (*entity.UserInfoToken, error)
	}

	User interface {
		GetUser(context.Context, int64) (*entity.User, error)
		GetUserByName(context.Context, string) (*entity.User, error)
		GetUsers(context.Context, entity.FilterUserInput) ([]*entity.User, error)
		ChangePassword(context.Context, entity.ChangePasswordInput) error
		UpdateRating(context.Context, int64, float32) error
		SetProfilePhoto(context.Context, int64, *multipart.FileHeader) error
	}

	Author interface {
		CreateAuthor(context.Context, entity.CreateAuthorInput) (*entity.Author, error)
		UpdateAuthor(context.Context, entity.UpdateAuthorInput) error
		GetAuthor(context.Context, int64) (*entity.Author, error)
		GetAuthors(context.Context) ([]*entity.Author, error)
		DeleteAuthor(context.Context, int64) error
	}

	Book interface {
		CreateBook(context.Context, entity.CreateBookInput) (*entity.Book, error)
		UpdateBook(context.Context, entity.UpdateBookInput) error
		GetBook(context.Context, int64) (*entity.Book, error)
		GetBooks(context.Context) ([]*entity.Book, error)
		DeleteBook(context.Context, int64) error
	}

	Export interface {
		GenerateExcelFile(context.Context) (*excelize.File, error)
		SaveToFile(*excelize.File) (string, error)
	}
	Command interface {
		UpdateCommands(context.Context, *multipart.FileHeader) error
	}
	Operation interface {
		CreateOperation(context.Context, entity.CreateOperationInput) (*entity.Operation, error)
		UpdateOperation(context.Context, entity.UpdateOperationInput) error
		DeleteOperation(context.Context, int64) error
	}
)
