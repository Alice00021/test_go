package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"test_go/internal/entity"
	"test_go/internal/repo"
	"test_go/pkg/auth"
	"test_go/pkg/jwt"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserService interface {
	Register(context.Context, *entity.UserInput) error
	Login(context.Context, string, string) (*TokenPair, error)
	RefreshTokens(context.Context, string) (*TokenPair, error)
	UpdateUser(context.Context, *entity.User) error
	GetByUserName(context.Context, string) (*entity.User, error)
	GetAllUsers(context.Context) ([]entity.User, error)
	ChangePassword(context.Context, string, string) error
	SetProfilePhoto(context.Context, uint, *multipart.FileHeader) error
	/* 	Logout(ctx context.Context, username string, password string) error */
}

type userUseCase struct {
	userRepo        repo.UserRepository
	jwtManager      *jwt.JWTManager
	storageBasePath string
}

func NewAuthService(userRepo repo.UserRepository, jwtManager *jwt.JWTManager, sbp string) UserService {
	return &userUseCase{
		userRepo:        userRepo,
		jwtManager:      jwtManager,
		storageBasePath: sbp,
	}
}

func (s *userUseCase) Register(ctx context.Context, inp *entity.UserInput) error {
	e := entity.NewUser(
		inp.Name, inp.Surname, inp.Username, inp.Password,
	)

	hashedPassword, err := auth.HashPassword(e.Password)
	if err != nil {
		return err
	}

	e.Password = hashedPassword
	return s.userRepo.Create(ctx, e)
}

func (s *userUseCase) Login(ctx context.Context, username string, password string) (*TokenPair, error) {
	user, err := s.userRepo.GetByUserName(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	acessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Username)

	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  acessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userUseCase) GetByUserName(ctx context.Context, username string) (*entity.User, error) {
	return s.userRepo.GetByUserName(ctx, username)
}

func (s *userUseCase) UpdateUser(ctx context.Context, user *entity.User) error {
	return s.userRepo.Update(ctx, user)
}

func (s *userUseCase) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *userUseCase) RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.jwtManager.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	user, err := s.userRepo.GetByUserName(ctx, claims.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	acessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Username)

	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  acessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *userUseCase) ChangePassword(ctx context.Context, username string, newpassword string) error {
	user, err := s.userRepo.GetByUserName(ctx, username)
	if err != nil {
		return errors.New("user not found")
	}
	hashedPassword, err := auth.HashPassword(newpassword)

	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.userRepo.Update(ctx, user)

}

func (s *userUseCase) SetProfilePhoto(ctx context.Context, id uint, file *multipart.FileHeader) error {
	allowedPhoto := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !allowedPhoto[ext] {
		return errors.New("invalid file type, only images (jpg, jpeg, png) are allowed")
	}

	if err := os.MkdirAll(s.storageBasePath, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory: %v", err)
	}

	user, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}
	filename := fmt.Sprintf("%s_profile%s", user.Username, ext)
	filePath := filepath.Join(s.storageBasePath, filename)

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}
	defer dst.Close()

	////////
	user.FilePath = &filePath
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user profile photo: %v", err)
	}
	return nil
}

/* func (s *userUseCase) Logout(){

}  */
