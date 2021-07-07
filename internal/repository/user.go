package repository

import (
	"context"

	"github.com/littlefut/go-template/pkg/errors"

	"github.com/littlefut/go-template/internal/user"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepo{db: db}
}

func (r *userRepo) FindByID(ctx context.Context, id int) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.NotFoundError, "user with id '%d' not found", id)
		}
		return nil, errors.New(errors.InternalError, err)
	}

	return &u, nil
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, "username = ?", username).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.NotFoundError, "user with username '%s' not found", username)
		}
		return nil, errors.New(errors.InternalError, err)
	}

	return &u, nil
}

func (r *userRepo) Save(ctx context.Context, u *user.User) error {
	_, err := r.FindByID(ctx, u.ID)
	if err == nil {
		return errors.New(errors.DuplicateError, "user with id: '%d' already exists", u.ID)
	}
	if errors.Is(err, errors.InternalError) {
		return err
	}

	err = r.db.WithContext(ctx).Create(u).Error
	if err != nil {
		return errors.New(errors.InternalError, err)
	}
	return nil
}

func (r *userRepo) Update(ctx context.Context, u *user.User) error {
	prevUser, err := r.FindByID(ctx, u.ID)
	if err != nil {
		return err
	}

	if u.Username != "" {
		prevUser.Username = u.Username
	}
	if !u.LastLogin.IsZero() {
		prevUser.LastLogin = u.LastLogin
	}

	err = r.db.WithContext(ctx).Save(prevUser).Error
	if err != nil {
		return errors.New(errors.InternalError, err)
	}
	return nil
}

func (r *userRepo) DeleteByID(ctx context.Context, id int) error {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).Delete(&user.User{}, "id = ?", id).Error
	if err != nil {
		return errors.New(errors.InternalError, err)
	}
	return nil
}
