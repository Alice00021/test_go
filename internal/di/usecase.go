package di

import (
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
	User   usecase.UserUseCase
	Book   usecase.BookUseCase
	Author usecase.AuthorService
}

func NewUseCase(
	t transactional.Transactional,
	repo *Repo,
	l logger.Interface,
	conf *config.Config,
	jwtManager *jwt.JWTManager,
) *UseCase {
	authorUc := author.New(t, repo.AuthorRepo, l)
	return &UseCase{
		Author: authorUc,
	}
}
