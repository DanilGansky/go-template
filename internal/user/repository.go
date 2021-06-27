package user

import (
	"context"
)

type Repository interface {
	FindByID(ctx context.Context, id int) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	Save(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	DeleteByID(ctx context.Context, id int) error
}
