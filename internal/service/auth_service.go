package service

import (
	"context"
	"errors"
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
	Register(ctx context.Context, user *entity.User) error 
	Login(ctx context.Context, username string, password string) (*TokenPair, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error)
    UpdateUser(ctx context.Context, user *entity.User, id uint) error
    GetByUserName(ctx context.Context, username string) (*entity.User, error)
    GetAllUsers(ctx context.Context) ([]entity.User, error)
	/* ChangePassword(ctx context.Context, username string, password string) error
	Logout(ctx context.Context, username string, password string) error */
}

type userService struct {
	userRepo repo.UserRepository
	jwtManager *jwt.JWTManager
}

func NewAuthService(userRepo repo.UserRepository, jwtManager *jwt.JWTManager) UserService {
	return &userService{
		userRepo: userRepo,
		jwtManager: jwtManager,
	}
}

func (s *userService) Register(ctx context.Context, user *entity.User) error {
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Create(ctx, user)
}

func (s *userService) Login(ctx context.Context, username string, password string) (*TokenPair, error) {
	user, err := s.userRepo.GetByUserName(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	acessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err!=nil{
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Username)

	if err!=nil{
		return nil, err
	}
	return &TokenPair{
		AccessToken: acessToken,
		RefreshToken :refreshToken,
	}, nil
}



func (s *userService) GetByUserName(ctx context.Context, username string)(*entity.User, error){
	return s.userRepo.GetByUserName(ctx, username)
}

func (s *userService) UpdateUser(ctx context.Context, user *entity.User, id uint) error{
	return s.userRepo.Update(ctx, user, id)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *userService) RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error){
	claims, err := s.jwtManager.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	user, err := s.userRepo.GetByUserName(ctx,claims.Username)
	if err!=nil{
		return nil, errors.New("user not found")
	}

	acessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Username)

	if err!=nil{
		return nil, err
	}

	return &TokenPair{
		AccessToken: acessToken,
		RefreshToken :newRefreshToken,
	}, nil	
}

/* func (s *userService) ChangePassword(){

}

func (s *userService) Logout(){

}  */