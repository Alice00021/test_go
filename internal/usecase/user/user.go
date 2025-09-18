package user

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
	"sync"
	"test_go/internal/entity"
	"test_go/internal/repo"
	"test_go/pkg/auth"
	"test_go/pkg/jwt"
	"test_go/pkg/logger"
	"test_go/pkg/transactional"

	"gopkg.in/gomail.v2"
)

type useCase struct {
	transactional.Transactional
	l               logger.Interface
	repo            repo.UserRepo
	jwtManager      *jwt.JWTManager
	storageBasePath string
	emailConfig     *entity.EmailConfig
	mtx             *sync.Mutex
}

func New(t transactional.Transactional,
	l logger.Interface,
	repo repo.UserRepo,
	jwtManager *jwt.JWTManager,
	sbp string,
	emailConfig *entity.EmailConfig,
	mtx *sync.Mutex,
) *useCase {
	return &useCase{
		Transactional:   t,
		l:               l,
		repo:            repo,
		jwtManager:      jwtManager,
		storageBasePath: sbp,
		emailConfig:     emailConfig,
		mtx:             mtx,
	}
}

func (uc *useCase) Register(ctx context.Context, inp *entity.CreateUserInput) (*entity.User, error) {
	op := "UserUseCase - Register"

	var user entity.User
	if err := uc.RunInTransaction(ctx, func(txCtx context.Context) error {
		_, err := uc.repo.GetByEmail(ctx, inp.Email)
		if err != nil {
			return entity.ErrEmailAlreadyUsed
		}

		verifyToken, err := generateVerifyToken()
		if err != nil {
			return entity.ErrGenerateVerifyToken
		}

		e := entity.NewUser(
			inp.Name, inp.Surname, inp.Username, inp.Password, inp.Email,
		)
		e.VerifyToken = &verifyToken

		e.Rating = 50

		hashedPassword, err := auth.HashPassword(e.Password)
		if err != nil {
			return err
		}

		e.Password = hashedPassword

		res, err := uc.repo.Create(txCtx, e)
		if err != nil {
			return fmt.Errorf("uc.repo.Create: %w", err)
		}

		if err := uc.sendVerificationEmail(e.Email, verifyToken); err != nil {
			return fmt.Errorf("uc.sendVerificationEmail: %w", err)
		}
		user = *res
		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s - uc.RunInTransaction: %w", op, err)
	}

	return &user, nil
}

func (uc *useCase) Login(ctx context.Context, username string, password string) (*entity.TokenPair, error) {
	op := "UserUseCase - Login"

	var tokenPair entity.TokenPair
	if err := uc.RunInTransaction(ctx, func(txCtx context.Context) error {
		user, err := uc.repo.GetByUserName(ctx, username)
		if err != nil {
			return fmt.Errorf("%s - uc.repo.GetByUserName: %w", op, err)
		}

		if !user.IsVerified {
			return fmt.Errorf("%s - %w", op, entity.ErrEmailNotVerified)
		}

		if !auth.CheckPasswordHash(password, user.Password) {
			return fmt.Errorf("%s - invalid credentials", op)
		}

		accessToken, err := uc.jwtManager.GenerateAccessToken(user.ID, user.Username)
		if err != nil {
			return fmt.Errorf("%s - uc.jwtManager.GenerateAccessToken: %w", op, err)
		}

		refreshToken, err := uc.jwtManager.GenerateRefreshToken(user.ID, user.Username)
		if err != nil {
			return fmt.Errorf("%s - uc.jwtManager.GenerateRefreshToken: %w", op, err)
		}
		tokenPair.AccessToken = accessToken
		tokenPair.RefreshToken = refreshToken

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%s - uc.RunInTransaction: %w", op, err)
	}

	return &tokenPair, nil
}

func (uc *useCase) GetUserByName(ctx context.Context, name string) (*entity.User, error) {
	user, err := uc.repo.GetByUserName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - GetUserByName - uc.repo.GetByUserName: %w", err)
	}

	return user, nil
}

func (uc *useCase) GetUser(ctx context.Context, id int64) (*entity.User, error) {
	user, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - GetUser - uc.repo.GetById: %w", err)
	}

	return user, nil
}

func (uc *useCase) GetUsers(ctx context.Context, filter entity.FilterUserInput) ([]*entity.User, error) {
	users, err := uc.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - GetUsers - uc.repo.GetAll: %w", err)
	}

	return users, nil
}

func (uc *useCase) UpdateUser(ctx context.Context, inp entity.UpdateUserInput) error {

	if err := uc.RunInTransaction(ctx, func(txCtx context.Context) error {
		e := &entity.User{
			Entity:   entity.Entity{ID: inp.ID},
			Name:     inp.Name,
			Surname:  inp.Surname,
			Username: inp.Username,
			Rating:   inp.Rating,
		}
		if err := uc.repo.Update(txCtx, e); err != nil {
			return fmt.Errorf("uc.repo.Update: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("%s - uc.RunInTransaction: %w", err)
	}

	return nil
}

func (uc *useCase) ChangePassword(ctx context.Context, id int64, oldPassword, newpassword, confirmPassword string) error {
	if newpassword != confirmPassword {
		return errors.New("newPassword and confirmPassword must be the same")
	}

	user, err := uc.repo.GetById(ctx, id)
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
	return uc.repo.Update(ctx, user)

}

func (uc *useCase) SetProfilePhoto(ctx context.Context, id int64, file *multipart.FileHeader) error {
	allowedPhoto := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !allowedPhoto[ext] {
		return errors.New("invalid file type, only images (jpg, jpeg, png) are allowed")
	}

	if err := os.MkdirAll(uc.storageBasePath, 0755); err != nil {
		return fmt.Errorf("failed to create upload directory: %v", err)
	}

	user, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}
	filename := fmt.Sprintf("%s_profile%s", user.Username, ext)
	filePath := filepath.Join(uc.storageBasePath, filename)

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
	if err := uc.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user profile photo: %v", err)
	}
	return nil
}

func (uc *useCase) VerifyEmail(ctx context.Context, token string) error {
	if err := uc.RunInTransaction(ctx, func(txCtx context.Context) error {
		user, err := uc.repo.GetByVerifyToken(ctx, token)
		if err != nil {
			return errors.New("invalid verification token")
		}

		user.IsVerified = true
		user.VerifyToken = nil

		if err := uc.repo.Update(txCtx, user); err != nil {
			return fmt.Errorf("uc.repo.Update: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("%s - uc.RunInTransaction: %w", err)
	}

	return nil
}

func (uc *useCase) RefreshTokens(ctx context.Context, refreshToken string) (*entity.TokenPair, error) {
	claims, err := uc.jwtManager.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}
	user, err := uc.repo.GetByUserName(ctx, claims.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	acessToken, err := uc.jwtManager.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := uc.jwtManager.GenerateRefreshToken(user.ID, user.Username)

	if err != nil {
		return nil, err
	}

	return &entity.TokenPair{
		AccessToken:  acessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func generateVerifyToken() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (us *useCase) sendVerificationEmail(email, token string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", us.emailConfig.SenderEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Email Verification")

	verificationLink := fmt.Sprintf("http://localhost:8080/auth/verify?token=%s", token)
	body := fmt.Sprintf("Please verify your email by clicking the following link: %s", verificationLink)
	message.SetBody("text/plain", body)

	d := gomail.NewDialer(us.emailConfig.SMTPHost, us.emailConfig.SMTPPort, us.emailConfig.SenderEmail, us.emailConfig.SenderPassword)

	return d.DialAndSend(message)
}
