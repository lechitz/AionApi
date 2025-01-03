package handlers

import (
	"github.com/lechitz/AionApi/ports/input"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

// USER STRUCTS

type User struct {
	UserService input.IUserService
	LoggerSugar *zap.SugaredLogger
}

type UserRequest struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        uint64         `json:"id"`
	Name      string         `json:"name"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type GetUserResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserResponse struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateUserResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LOGIN STRUCTS

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
