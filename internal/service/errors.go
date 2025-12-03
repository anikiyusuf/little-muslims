package service 

import (
	"errors"
)


var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
    ErrSetVerificationCode = errors.New("failed to set verification code in cache")
    ErrInvalidCode		   = errors.New("invalid verification code")
	ErrUserNotFound		   = errors.New("user not found")
    ErrGenerateToken	   = errors.New("failed to generate access token")
)