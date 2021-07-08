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
		return errors.New(errors.InternalError, err.Error())
	}

	user := User{
		Username: dto.Username,
		Password: hashedPassword,
		JoinedAt: time.Now(),
	}

	if err = s.repo.Save(ctx, &user); err != nil {
		return errors.New(errors.InternalError, err.Error())
	}
	return nil
}

func (s *service) FindByID(ctx context.Context, id int) (*DTO, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New(errors.NotFoundError, err.Error())
	}
	return MakeDTO(user), nil
}

func (s *service) FindByUsername(ctx context.Context, username string) (*DTO, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New(errors.NotFoundError, err.Error())
	}
	return MakeDTO(user), nil
}

func (s *service) SetUsername(ctx context.Context, id int, dto *UpdateDTO) error {
	if dto.Username == "" {
		return ErrEmptyUsername
	}

	err := s.repo.Update(ctx, &User{ID: id, Username: dto.Username})
	if err != nil {
		return errors.New(errors.InternalError, err.Error())
	}
	return nil
}

func (s *service) SetLastLogin(ctx context.Context, id int, lastLogin time.Time) error {
	if lastLogin.IsZero() {
		return ErrEmptyLastLogin
	}

	err := s.repo.Update(ctx, &User{ID: id, LastLogin: lastLogin})
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	err := s.repo.DeleteByID(ctx, id)
	if err != nil {
		return errors.New(errors.NotFoundError, err.Error())
	}
	return nil
}
