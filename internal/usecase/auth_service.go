package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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

	"gopkg.in/gomail.v2"
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
	ChangePassword(context.Context, uint, string, string, string) error
	SetProfilePhoto(context.Context, uint, *multipart.FileHeader) error
	VerifyEmail(context.Context, string) error
	/* 	Logout(ctx context.Context, username string, password string) error */
}

type userUseCase struct {
	userRepo        repo.UserRepository
	jwtManager      *jwt.JWTManager
	storageBasePath string
	emailConfig     *entity.EmailConfig
}

func NewAuthService(userRepo repo.UserRepository, jwtManager *jwt.JWTManager, sbp string, emailConfig *entity.EmailConfig) UserService {
	return &userUseCase{
		userRepo:        userRepo,
		jwtManager:      jwtManager,
		storageBasePath: sbp,
		emailConfig:     emailConfig,
	}
}

func (s *userUseCase) Register(ctx context.Context, inp *entity.UserInput) error {
	existingUserByEmail, err := s.userRepo.GetByEmail(ctx, inp.Email)
	if err == nil && existingUserByEmail != nil {
		return errors.New("email already exists")
	}

	verifyToken, err := generateVerifyToken()
	if err != nil {
		return fmt.Errorf("failed to generate verify token: %v", err)
	}

	e := entity.NewUser(
		inp.Name, inp.Surname, inp.Username, inp.Password, inp.Email,
	)
	e.VerifyToken = &verifyToken

	hashedPassword, err := auth.HashPassword(e.Password)
	if err != nil {
		return err
	}

	e.Password = hashedPassword

	if err := s.userRepo.Create(ctx, e); err != nil {
		return err
	}

	if err := s.sendVerificationEmail(e.Email, verifyToken); err != nil {
		return fmt.Errorf("failed to send verification email: %v", err)
	}

	return nil
}

func (s *userUseCase) Login(ctx context.Context, username string, password string) (*TokenPair, error) {
	user, err := s.userRepo.GetByUserName(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !user.IsVerified {
		return nil, entity.ErrEmailNotVerified
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

func (s *userUseCase) ChangePassword(ctx context.Context, id uint, oldPassword, newpassword, confirmPassword string) error {
	if newpassword != confirmPassword {
		return errors.New("newPassword and confirmPassword must be the same")
	}

	user, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	if !auth.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("error change password")
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

func (s *userUseCase) VerifyEmail(ctx context.Context, token string) error {
	user, err := s.userRepo.GetByVerifyToken(ctx, token)
	if err != nil {
		return errors.New("invalid verification token")
	}

	user.IsVerified = true
	user.VerifyToken = nil
	return s.userRepo.Update(ctx, user)
}

func generateVerifyToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *userUseCase) sendVerificationEmail(email, token string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", s.emailConfig.SenderEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Email Verification")

	verificationLink := fmt.Sprintf("http://localhost:8080/auth/verify?token=%s", token)
	body := fmt.Sprintf("Please verify your email by clicking the following link: %s", verificationLink)
	message.SetBody("text/plain", body)

	d := gomail.NewDialer(s.emailConfig.SMTPHost, s.emailConfig.SMTPPort, s.emailConfig.SenderEmail, s.emailConfig.SenderPassword)

	return d.DialAndSend(message)
}

/* func (s *userUseCase) Logout(){

}  */
