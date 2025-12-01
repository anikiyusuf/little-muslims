package service 

import (
	"errors"
)


var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
    ErrSetVerificationCode = errors.New("failed to set verification code in cache")
)