package handlers

import (
	"github.com/lechitz/AionApi/ports/input"
	"go.uber.org/zap"
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
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
}

type GetUserResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserResponse struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

// LOGIN STRUCTS

type Login struct {
	LoginService input.ILoginService
	LoggerSugar  *zap.SugaredLogger
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
