package user

import "time"

func MakeDTO(u *User) *DTO {
	return &DTO{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		JoinedAt:  u.JoinedAt.Format(time.RFC822),
		LastLogin: u.LastLogin.Format(time.RFC822),
	}
}
