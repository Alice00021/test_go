package entity

import "errors"

var (
	ErrAccessDenied        = errors.New("access denied")
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailNotVerified    = errors.New("email not verified")
	ErrEmailAlreadyUsed    = errors.New("email already used")
	ErrAuthorNotFound      = errors.New("author not found")
	ErrBookNotFound        = errors.New("book not found")
	ErrGenerateVerifyToken = errors.New("failed to generate verify token")
)
