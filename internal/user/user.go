package user

import (
	"fmt"
	"time"
)

type User struct {
	ID       int    `gorm:"primaryKey;column:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`

	JoinedAt  time.Time `gorm:"column:joined_at"`
	LastLogin time.Time `gorm:"column:last_login"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) String() string {
	return fmt.Sprintf("user: id '%d'; username: '%s'", u.ID, u.Username)
}
