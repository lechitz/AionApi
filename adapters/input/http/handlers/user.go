package handlers

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/utils"
	"github.com/lechitz/AionApi/ports/input"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	ErrorToCreateUser = "error to create and process the request"
)

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

func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var userRequest UserRequest
	json.NewDecoder(r.Body).Decode(&userRequest)

	var userDomain domain.UserDomain
	copier.Copy(&userDomain, &userRequest)

	userDomain, err := u.UserService.CreateUser(contextControl, userDomain)
	if err != nil {
		u.LoggerSugar.Errorw(ErrorToCreateUser, "error", err.Error())
		response := utils.ObjectResponse(ErrorToCreateUser, err.Error())
		utils.ResponseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var userResponse UserResponse
	copier.Copy(&userResponse, &userDomain)
	response := utils.ObjectResponse(userResponse, "User created with success")
	utils.ResponseReturn(w, http.StatusCreated, response.Bytes())
}
