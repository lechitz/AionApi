package dto

import (
	"time"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ID       uint64 `json:"user_id"`
}

type CreateUserResponse struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       uint64 `json:"user_id"`
}

type GetUserResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	ID        uint64    `json:"user_id"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	ID       uint64  `json:"user_id"`
}

type UpdateUserResponse struct {
	UpdatedAt time.Time `json:"updated_at"`
	Name      *string   `json:"name"`
	Username  *string   `json:"username"`
	Email     *string   `json:"email"`
	ID        uint64    `json:"user_id"`
}

type UpdatePasswordUserRequest struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}
