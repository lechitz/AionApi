package dto

// CreateUserResponse is the HTTP output for user creation.
type CreateUserResponse struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
