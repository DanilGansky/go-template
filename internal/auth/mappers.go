package auth

import (
	"time"

	"github.com/littlefut/go-template/internal/user"
)

func MakeLoggedInDTO(dto *user.DTO, token string, lastLogin time.Time) *LoggedInDTO {
	return &LoggedInDTO{
		Username:  dto.Username,
		LastLogin: lastLogin.Format(time.RFC822),
		JoinedAt:  dto.JoinedAt,
		Token:     token,
	}
}
