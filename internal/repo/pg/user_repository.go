package pg

import (
	"context"
	"gorm.io/gorm"
	"test_go/internal/entity"
	"test_go/internal/repo"
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

func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepo) GetByUserName(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetAll(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) GetById(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByVerifyToken(ctx context.Context, token string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Where("verify_token = ?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
