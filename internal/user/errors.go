package user

import "github.com/littlefut/go-template/pkg/errors"

var (
	ErrEmptyUsername  = errors.New(errors.ValidationError, "username cannot be empty")
	ErrEmptyPassword  = errors.New(errors.ValidationError, "password cannot be empty")
	ErrEmptyLastLogin = errors.New(errors.ValidationError, "lastLogin cannot be empty")
)
