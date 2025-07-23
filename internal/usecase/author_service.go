package usecase

import (
	"context"

	"test_go/internal/entity"
	"test_go/internal/repo"
)

type AuthorService interface {
	CreateAuthor(ctx context.Context, name string, gender bool) (*entity.Author, error)
	UpdateAuthor(ctx context.Context, name string, gender bool, id uint) (*entity.Author, error)
	DeleteAuthor(ctx context.Context, id uint) error
	GetByIDAuthor(ctx context.Context, id uint) (*entity.Author, error)
	GetAllAuthors(ctx context.Context) ([]entity.Author, error)
}

type authorService struct {
	repo repo.AuthorRepository
}

func NewAuthorService(repo repo.AuthorRepository) AuthorService {
	return &authorService{repo: repo}
}

func (s *authorService) CreateAuthor(ctx context.Context, name string, gender bool) (*entity.Author, error) {
	author := &entity.Author{Name: name, Gender: gender}
	if err := s.repo.Create(ctx, author); err != nil {
		return nil, err
	}
	return author, nil
}

func (s *authorService) UpdateAuthor(ctx context.Context, name string, gender bool, id uint) (*entity.Author, error) {
	author, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	author.Name = name
	author.Gender = gender

	if err := s.repo.Update(ctx, author, id); err != nil {
		return nil, err
	}
	return author, nil
}

func (s *authorService) DeleteAuthor(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *authorService) GetByIDAuthor(ctx context.Context, id uint) (*entity.Author, error) {
	book, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *authorService) GetAllAuthors(ctx context.Context) ([]entity.Author, error) {
	author, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return author, nil
}
