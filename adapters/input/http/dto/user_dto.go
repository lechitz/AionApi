package dto

import (
	"time"
)

type CreateUserRequest struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type GetUserResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateUserRequest struct {
	ID       uint64  `json:"id"`
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

type UpdateUserResponse struct {
	ID        uint64    `json:"id"`
	Name      *string   `json:"name"`
	Username  *string   `json:"username"`
	Email     *string   `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdatePasswordRequest struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}
