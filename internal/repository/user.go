package repository

import (
	"context"
	"fmt"

	"github.com/littlefut/go-template/internal/user"
	"github.com/pkg/errors"
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
		if errors.Cause(err) == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(user.ErrNotFound, fmt.Sprintf("user with id '%d' not found", id))
		}
		return nil, errors.Wrap(err, "internal error")
	}

	return &u, nil
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, "username = ?", username).Error; err != nil {
		if errors.Cause(err) == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(user.ErrNotFound, fmt.Sprintf("user with username '%s' not found", username))
		}
		return nil, errors.Wrap(err, "internal error")
	}

	return &u, nil
}

func (r *userRepo) Save(ctx context.Context, u *user.User) error {
	_, err := r.FindByID(ctx, u.ID)
	if err == nil {
		return fmt.Errorf("user with id: '%d' already exists", u.ID)
	}
	if errors.Cause(err) == user.ErrInternalError {
		return errors.Wrap(err, "failed to create user")
	}

	err = r.db.WithContext(ctx).Create(u).Error
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}

func (r *userRepo) Update(ctx context.Context, u *user.User) error {
	user, err := r.FindByID(ctx, u.ID)
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	if u.Username != "" {
		user.Username = u.Username
	}
	if !u.LastLogin.IsZero() {
		user.LastLogin = u.LastLogin
	}

	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return errors.Wrap(err, "failed to create user")
	}
	return nil
}

func (r *userRepo) DeleteByID(ctx context.Context, id int) error {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	if err := r.db.Delete(&user.User{}, "id = ?", id).Error; err != nil {
		return errors.Wrap(err, "internal error")
	}
	return nil
}
