package user

type RegisterDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateDTO struct {
	Username string `json:"username"`
}

type CredentialsDTO struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	JoinedAt  string `json:"joined_at"`
	LastLogin string `json:"last_login"`
}

type DTO struct {
	Username  string `json:"username"`
	JoinedAt  string `json:"joined_at"`
	LastLogin string `json:"last_login"`
}
