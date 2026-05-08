package response

type RegisterResponse struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}
