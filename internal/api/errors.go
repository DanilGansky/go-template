package api

import "github.com/littlefut/go-template/pkg/errors"

var (
	ErrInvalidContext    = errors.New(errors.InternalError, "invalid context")
	ErrTokenDoesNotValid = errors.New(errors.ValidationError, "token does not valid")
)
