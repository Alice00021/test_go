package request

import "test_go/internal/entity"

type AuthenticateRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"oldPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
}

func (req *ChangePasswordRequest) ToEntity() entity.ChangePasswordInput {
	return entity.ChangePasswordInput{
		OldPassword:     req.OldPassword,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}
}

type VerifyEmailRequest struct {
	Token string `form:"token" validate:"required"`
}
type UpdateRatingRequest struct {
	Rating float32 `json:"rating" validate:"required" min:"0" max:"100"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

func (req *CreateUserRequest) ToEntity() entity.CreateUserInput {
	return entity.CreateUserInput{
		Name:     req.Name,
		Surname:  req.Surname,
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}
}
