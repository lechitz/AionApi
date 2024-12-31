package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/badoux/checkmail"
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/utils"
	"github.com/lechitz/AionApi/ports/input"
	"go.uber.org/zap"
	"net/http"
	"strings"
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
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		u.LoggerSugar.Errorw(ErrorToCreateUser, "error", err.Error())
		response := utils.ObjectResponse(ErrorToCreateUser, err.Error())
		utils.ResponseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	if err := userRequest.prepareUser("register"); err != nil {
		u.LoggerSugar.Errorw(ErrorToCreateUser, "error", err.Error())
		response := utils.ObjectResponse(ErrorToCreateUser, err.Error())
		utils.ResponseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

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

func (userRequest *UserRequest) prepareUser(step string) error {
	if err := userRequest.validate(step); err != nil {
		return err
	}

	if err := userRequest.format(step); err != nil {
		return err
	}

	return nil
}

func (userRequest *UserRequest) validate(step string) error {
	if userRequest.Name == "" {
		return errors.New("name is required")
	}
	if userRequest.Username == "" {
		return errors.New("username is required")
	}
	if userRequest.Email == "" {
		return errors.New("email is required")
	}

	if err := checkmail.ValidateFormat(userRequest.Email); err != nil {
		return errors.New("invalid email")
	}

	if step == "register" && userRequest.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (userRequest *UserRequest) format(step string) error {
	userRequest.Name = strings.TrimSpace(userRequest.Name)
	userRequest.Username = strings.TrimSpace(userRequest.Username)
	userRequest.Email = strings.TrimSpace(userRequest.Email)

	if step == "register" {
		hashedPassword, err := utils.Hash(userRequest.Password)
		if err != nil {
			return err
		}

		userRequest.Password = string(hashedPassword)
	}
	return nil
}
