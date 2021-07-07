package auth

import "github.com/littlefut/go-template/pkg/errors"

var (
	ErrEmptyUsername   = errors.New(errors.ValidationError, "username cannot be empty")
	ErrEmptyPassword   = errors.New(errors.ValidationError, "password cannot be empty")
	ErrInvalidPassword = errors.New(errors.Unauthorized, "invalid password")
)
