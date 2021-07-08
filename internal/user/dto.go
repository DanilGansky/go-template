package user

type RegisterDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateDTO struct {
	Username string `json:"username"`
}

type DTO struct {
	ID        int    `json:"id,omitempty"`
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	JoinedAt  string `json:"joined_at"`
	LastLogin string `json:"last_login"`
}
