package auth

import "errors"

var (
	ErrTokenDoesNotValid = errors.New("token does not valid")
	ErrValidation        = errors.New("validation error")
	ErrInvalidPassword   = errors.New("passwords are mismatch")
)
