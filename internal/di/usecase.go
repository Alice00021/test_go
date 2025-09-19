package di

import (
	"sync"
	"test_go/config"
	"test_go/internal/usecase"
	"test_go/internal/usecase/author"
	"test_go/internal/usecase/book"
	"test_go/internal/usecase/export"
	"test_go/internal/usecase/user"
	"test_go/pkg/jwt"
	"test_go/pkg/logger"
	"test_go/pkg/transactional"
)

type UseCase struct {
	User   usecase.User
	Book   usecase.Book
	Author usecase.Author
	Export usecase.Export
}

func NewUseCase(
	t transactional.Transactional,
	repo *Repo,
	l logger.Interface,
	conf *config.Config,
	jwtManager *jwt.JWTManager,
) *UseCase {
	txMtx := &sync.Mutex{}
	userUc := user.New(t, l, repo.UserRepo, jwtManager, conf.LocalFileStorage.BasePath, &conf.EmailConfig, txMtx)
	authorUc := author.New(t, repo.AuthorRepo, l)
	bookUc := book.New(t, repo.BookRepo, l)
	exportUc := export.New(authorUc, bookUc, l, conf.LocalFileStorage.ExportPath)
	return &UseCase{
		Author: authorUc,
		Book:   bookUc,
		User:   userUc,
		Export: exportUc,
	}
}
