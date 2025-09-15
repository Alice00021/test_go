package entity

import "errors"

var (
	ErrAccessDenied     = errors.New("access denied")
	ErrEmailNotVerified = errors.New("email not verified")
)
