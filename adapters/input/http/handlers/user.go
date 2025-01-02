package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/badoux/checkmail"
	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/utils"
	"net/http"
	"strconv"
	"strings"
)

const (
	ErrorToDecodeUserRequest = "error to decode user request"
	ErrorToPrepareUser       = "error to prepare user"
	ErrorToCreateUser        = "error to create user"
	ErrorToGetUser           = "error to get user"
	ErrorToGetUsers          = "error to get users"

	ErrorToParseUser = "error to parse user"

	SuccessToCreateUser = "user created successfully"
	SucessToGetUser     = "user get successfully"
	SucessToGetUsers    = "users get successfully"

	MissingUserIDParameter = "missing user ID parameter"

	NameIsRequired     = "name is required"
	UsernameIsRequired = "username is required"
	EmailIsRequired    = "email is required"
	PasswordIsRequired = "password is required"
	UserIDIsRequired   = "user ID is required"

	InvalidEmail = "invalid email"
)

func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var userRequest UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		u.LoggerSugar.Errorw(ErrorToDecodeUserRequest, "error", err.Error())
		response := utils.ObjectResponse(ErrorToDecodeUserRequest, err.Error())
		utils.ResponseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	if err := userRequest.prepareUser("register"); err != nil {
		u.LoggerSugar.Errorw(ErrorToPrepareUser, "error", err.Error())
		response := utils.ObjectResponse(ErrorToPrepareUser, err.Error())
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
	response := utils.ObjectResponse(userResponse, SuccessToCreateUser)
	utils.ResponseReturn(w, http.StatusCreated, response.Bytes())
}

func (u *User) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	usersDomain, err := u.UserService.GetAllUsers(contextControl)
	if err != nil {
		u.LoggerSugar.Errorw(ErrorToGetUsers, "error", err.Error())
		response := utils.ObjectResponse(ErrorToGetUsers, err.Error())
		utils.ResponseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var getUsersResponse []GetUserResponse
	copier.Copy(&getUsersResponse, &usersDomain)
	response := utils.ObjectResponse(getUsersResponse, SucessToGetUsers)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (u *User) GetUserByID(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	userIDParam := chi.URLParam(r, "id")

	if userIDParam == "" {
		u.LoggerSugar.Errorw(MissingUserIDParameter)
		response := utils.ObjectResponse(MissingUserIDParameter, UserIDIsRequired)
		utils.ResponseReturn(w, http.StatusBadRequest, response.Bytes())
		return
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		u.LoggerSugar.Errorw(ErrorToParseUser, "error", err.Error())
		response := utils.ObjectResponse(ErrorToParseUser, err.Error())
		utils.ResponseReturn(w, http.StatusBadRequest, response.Bytes())
		return
	}

	userDomain, err := u.UserService.GetUserByID(contextControl, userID)
	if err != nil {
		u.LoggerSugar.Errorw(ErrorToGetUser, "error", err.Error())
		response := utils.ObjectResponse(ErrorToGetUser, err.Error())
		utils.ResponseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var getUserResponse []GetUserResponse
	copier.Copy(&getUserResponse, &userDomain)
	response := utils.ObjectResponse(getUserResponse, SucessToGetUser)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
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
		return errors.New(NameIsRequired)
	}
	if userRequest.Username == "" {
		return errors.New(UsernameIsRequired)
	}
	if userRequest.Email == "" {
		return errors.New(EmailIsRequired)
	}

	if err := checkmail.ValidateFormat(userRequest.Email); err != nil {
		return errors.New(InvalidEmail)
	}

	if step == "register" && userRequest.Password == "" {
		return errors.New(PasswordIsRequired)
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
