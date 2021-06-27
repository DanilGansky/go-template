package user

import "errors"

var (
	ErrInternalError = errors.New("internal error")
	ErrNotFound      = errors.New("failed to find user")
	ErrValidation    = errors.New("validation error")
)
