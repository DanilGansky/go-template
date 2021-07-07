package auth

import (
	"context"
	"time"

	"github.com/littlefut/go-template/pkg/errors"

	"github.com/littlefut/go-template/internal/hash"
	"github.com/littlefut/go-template/internal/user"
)

type Service interface {
	Login(ctx context.Context, dto *LoginDTO) (*DTO, error)
}

type service struct {
	hashSvc  hash.Service
	tokenSvc hash.TokenService
	userSvc  user.Service
}

func NewService(hashSvc hash.Service, tokenSvc hash.TokenService, userSvc user.Service) Service {
	return &service{
		hashSvc:  hashSvc,
		tokenSvc: tokenSvc,
		userSvc:  userSvc,
	}
}

func (s *service) Login(ctx context.Context, dto *LoginDTO) (*DTO, error) {
	if dto.Username == "" {
		return nil, ErrEmptyUsername
	}
	if dto.Password == "" {
		return nil, ErrEmptyPassword
	}

	credentials, err := s.userSvc.FindCredentialsByUsername(ctx, dto.Username)
	if err != nil {
		return nil, err
	}
	if !s.hashSvc.Compare(credentials.Password, dto.Password) {
		return nil, ErrInvalidPassword
	}

	token, err := s.tokenSvc.Generate(credentials.ID)
	if err != nil {
		return nil, errors.New(errors.InternalError, err)
	}

	err = s.userSvc.SetLastLogin(ctx, credentials.ID, time.Now())
	if err != nil {
		return nil, err
	}

	return &DTO{
		Username:  credentials.Username,
		LastLogin: time.Now().Format("02 Jan 06 15:04 MST"),
		JoinedAt:  credentials.JoinedAt,
		Token:     token,
	}, nil
}
