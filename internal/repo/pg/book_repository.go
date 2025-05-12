package pg

import (
    "context"
    "test_go/internal/entity"
    "test_go/internal/repo"
    "gorm.io/gorm"
)

type BookRepo struct {
    db *gorm.DB
}

// NewBookRepo — конструктор для создания экземпляра BookRepo
func NewBookRepo(db *gorm.DB) repo.BookRepository {
    return &BookRepo{db: db}
}

func (r *BookRepo) Create(ctx context.Context, book *entity.Book) error {
    return r.db.WithContext(ctx).Create(book).Error
}

func (r *BookRepo) Update(ctx context.Context, book *entity.Book) error {
    return r.db.WithContext(ctx).Save(book).Error
}

func (r *BookRepo) Delete(ctx context.Context, id uint) error {
    result := r.db.WithContext(ctx).Delete(&entity.Book{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    return nil
}

func (r *BookRepo) GetByID(ctx context.Context, id uint) (*entity.Book, error) {
    var book entity.Book
    if err := r.db.WithContext(ctx).Preload("Author").First(&book, id).Error; err != nil {
        return nil, err
    }
    return &book, nil
}

func (r *BookRepo) GetAll(ctx context.Context) ([]entity.Book, error) {
    var books []entity.Book
    if err := r.db.WithContext(ctx).Preload("Author").Find(&books).Error; err != nil {
        return nil, err
    }
    return books, nil
}