package pg

import (
	"context"
	"test_go/internal/entity"
	"test_go/internal/repo"

	"gorm.io/gorm"
)

// TranslationRepo -.
type AuthorRepo struct {
	db *gorm.DB
}


// NewAuthorPostgres — конструктор для создания экземпляра AuthorRepo
func NewAuthorRepo(db *gorm.DB) repo.AuthorRepository {
	return &AuthorRepo{db: db}
}

func (r *AuthorRepo) Create(ctx context.Context, author *entity.Author) error {
	return r.db.WithContext(ctx).Create(author).Error
}

func (r *AuthorRepo) Update(ctx context.Context, author *entity.Author) error {
	return r.db.WithContext(ctx).Save(author).Error
}
func (r *AuthorRepo) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entity.Author{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil

}

func (r *AuthorRepo) GetAll(ctx context.Context) ([]entity.Author, error) {
	var authors []entity.Author
    if err := r.db.WithContext(ctx).Find(&authors).Error; err != nil {
        return nil, err
    }
    return authors, nil
}

func (r *AuthorRepo) GetByID(ctx context.Context, id uint) (*entity.Author, error) {
	var author entity.Author
    if err := r.db.WithContext(ctx).First(&author, id).Error; err != nil {
        return nil, err
    }
    return &author, nil
}

