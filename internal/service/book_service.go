package service

import (
    "context"
    "test_go/internal/entity"
    "test_go/internal/repo"
)

type BookService interface {
    CreateBook(ctx context.Context, title string, authorID uint) (*entity.Book, error)
    
}

type bookService struct {
    repo repo.BookRepository
}

func NewBookService(repo repo.BookRepository) BookService {
    return &bookService{repo: repo}
}

func (s *bookService) CreateBook(ctx context.Context, title string, authorID uint) (*entity.Book, error) {
    book := &entity.Book{Title: title, AuthorID: authorID}
    if err := s.repo.Create(ctx, book); err != nil {
        return nil, err
    }
    return book, nil
}