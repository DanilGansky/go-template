package auth

import (
	"context"
	"time"

	"github.com/littlefut/go-template/pkg/errors"

	"github.com/littlefut/go-template/internal/hash"
	"github.com/littlefut/go-template/internal/user"
)

type Service interface {
	Login(ctx context.Context, dto *LoginDTO) (*LoggedInDTO, error)
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

func (s *service) Login(ctx context.Context, dto *LoginDTO) (*LoggedInDTO, error) {
	if dto.Username == "" {
		return nil, ErrEmptyUsername
	}
	if dto.Password == "" {
		return nil, ErrEmptyPassword
	}

	userDTO, err := s.userSvc.FindByUsername(ctx, dto.Username)
	if err != nil {
		return nil, err
	}
	if !s.hashSvc.Compare(userDTO.Password, dto.Password) {
		return nil, ErrInvalidPassword
	}

	token, err := s.tokenSvc.Generate(userDTO.ID)
	if err != nil {
		return nil, errors.New(errors.InternalError, err)
	}

	lastLogin := time.Now()
	if err = s.userSvc.SetLastLogin(ctx, userDTO.ID, lastLogin); err != nil {
		return nil, err
	}
	return MakeLoggedInDTO(userDTO, token, lastLogin), nil
}
