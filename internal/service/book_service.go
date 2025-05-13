package service

import (
    "context"
    "test_go/internal/entity"
    "test_go/internal/repo"
)

type BookService interface {
    CreateBook(ctx context.Context, title string, authorID uint) (*entity.Book, error)
    UpdateBook(ctx context.Context, title string, authorID uint) (*entity.Book, error)
    DeleteBook(ctx context.Context, id uint)   error
    GetByIDBook(ctx context.Context, id uint) (*entity.Book, error)
    GetAllBooks(ctx context.Context) ([]entity.Book, error) 
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

func (s *bookService) UpdateBook(ctx context.Context, title string, authorID uint) (*entity.Book, error) {
    book := &entity.Book{Title: title, AuthorID: authorID}
    if err := s.repo.Update(ctx, book); err != nil {
        return nil, err
    }
    return book, nil
}


func (s *bookService) DeleteBook(ctx context.Context, id uint)  error {
    if err := s.repo.Delete(ctx, id); err != nil {
        return err
    }
    return  nil
}

func (s *bookService) GetByIDBook(ctx context.Context, id uint)  (*entity.Book, error) {
    book, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    return  book, nil
}

func (s *bookService) GetAllBooks(ctx context.Context) ([]entity.Book, error){
    books, err := s.repo.GetAll(ctx)
    if err!=nil{
        return nil, err
    }
    return books, nil
}
