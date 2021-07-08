package user

import (
	"context"
	"time"

	"github.com/littlefut/go-template/pkg/errors"

	"github.com/littlefut/go-template/internal/hash"
)

type Service interface {
	Register(ctx context.Context, dto *RegisterDTO) error

	FindByID(ctx context.Context, id int) (*DTO, error)
	FindByUsername(ctx context.Context, username string) (*DTO, error)

	SetUsername(ctx context.Context, id int, dto *UpdateDTO) error
	SetLastLogin(ctx context.Context, id int, lastLogin time.Time) error

	Delete(ctx context.Context, id int) error
}

type service struct {
	hashSvc hash.Service
	repo    Repository
}

func NewService(hashSvc hash.Service, repo Repository) Service {
	return &service{hashSvc: hashSvc, repo: repo}
}

func (s *service) Register(ctx context.Context, dto *RegisterDTO) error {
	if dto.Username == "" {
		return ErrEmptyUsername
	}
	if dto.Password == "" {
		return ErrEmptyPassword
	}

	hashedPassword, err := s.hashSvc.Encrypt(dto.Password)
	if err != nil {
		return errors.New(errors.InternalError, err)
	}

	user := User{
		Username: dto.Username,
		Password: hashedPassword,
		JoinedAt: time.Now(),
	}
	return s.repo.Save(ctx, &user)
}

func (s *service) FindByID(ctx context.Context, id int) (*DTO, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return MakeDTO(user), nil
}

func (s *service) FindByUsername(ctx context.Context, username string) (*DTO, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return MakeDTO(user), nil
}

func (s *service) SetUsername(ctx context.Context, id int, dto *UpdateDTO) error {
	if dto.Username == "" {
		return ErrEmptyUsername
	}
	return s.repo.Update(ctx, &User{ID: id, Username: dto.Username})
}

func (s *service) SetLastLogin(ctx context.Context, id int, lastLogin time.Time) error {
	if lastLogin.IsZero() {
		return ErrEmptyLastLogin
	}
	return s.repo.Update(ctx, &User{ID: id, LastLogin: lastLogin})
}

func (s *service) Delete(ctx context.Context, id int) error {
	return s.repo.DeleteByID(ctx, id)
}
