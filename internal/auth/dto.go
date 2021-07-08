package auth

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoggedInDTO struct {
	Username  string `json:"username"`
	LastLogin string `json:"last_login"`
	JoinedAt  string `json:"joined_at"`
	Token     string `json:"token"`
}
