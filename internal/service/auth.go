package service

import (
	"context"
	"errors"
	"test_go/internal/entity"
	"test_go/internal/repo"
	"test_go/pkg/auth"
	"test_go/pkg/jwt"
)

type AuthService struct {
	userRepo repo.UserRepository
}

func NewAuthService(userRepo repo.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, user *entity.User) error {
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetByUserName(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}