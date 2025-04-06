package dto

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Username string `json:"username"`
}

type LogoutUserRequest struct {
	UserID uint64 `json:"id"`
	Token  string `json:"token"`
}
