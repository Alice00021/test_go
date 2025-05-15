package pg

import (
    "context"
    "test_go/internal/entity"
    "test_go/internal/repo"
    "gorm.io/gorm"
)

type UserRepo struct {
    db *gorm.DB
}

// NewUserRepo — конструктор для создания экземпляра UserRepo
func NewUserRepo(db *gorm.DB) repo.UserRepository {
    return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) Update(ctx context.Context, user *entity.User, id uint) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepo) Delete(ctx context.Context, id uint) error {
    result := r.db.WithContext(ctx).Delete(&entity.User{}, id)
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    return nil
}

func (r *UserRepo) GetByUserName(ctx context.Context, username string) (*entity.User, error) {
    var user entity.User
    if err := r.db.WithContext(ctx).Preload("Author").First(&user, username).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepo) GetAll(ctx context.Context) ([]entity.User, error) {
    var users []entity.User
    if err := r.db.WithContext(ctx).Preload("Author").Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}