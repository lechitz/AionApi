package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/badoux/checkmail"
	"github.com/go-chi/chi/v5"
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/adapters/middlewares"
	"github.com/lechitz/AionApi/adapters/security"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/pkg/utils"
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
	ErrorToExtractUserID     = "error to extract user ID"
	ErrorToUpdateUser        = "error to update user"

	ErrUserPermissionDenied = "user permission denied"
	ErrorToParseUser        = "error to parse user"

	SuccessToCreateUser = "user created successfully"
	SuccessToGetUser    = "user get successfully"
	SuccessToGetUsers   = "users get successfully"
	SuccessToUpdateUser = "user updated successfully"
	SuccessToDeleteUser = "user deleted successfully"

	MissingUserIDParameter = "missing user ID parameter"

	NameIsRequired     = "name is required"
	UsernameIsRequired = "username is required"
	EmailIsRequired    = "email is required"
	PasswordIsRequired = "password is required"
	UserIDIsRequired   = "user ID is required"

	InvalidEmail = "invalid email"

	ErrorToDecodeLoginRequest = "error to decode login request"
	ErrorToGetUserByUsername  = "error to get user by username"
	ErrorToVerifyPassword     = "error to verify password"
	ErrorToCreateToken        = "error to create token"

	SuccessToLogin = "success to login"
)

func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var userRequest UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, ErrorToDecodeUserRequest, err)
		return
	}

	if err := userRequest.prepareUser("register"); err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, ErrorToPrepareUser, err)
		return
	}

	var userDomain domain.UserDomain
	copier.Copy(&userDomain, &userRequest)

	userDomain, err := u.UserService.CreateUser(contextControl, userDomain)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, ErrorToCreateUser, err)
		return
	}

	var createUserResponse CreateUserResponse
	copier.Copy(&createUserResponse, &userDomain)
	response := utils.ObjectResponse(createUserResponse, SuccessToCreateUser)
	utils.ResponseReturn(w, http.StatusCreated, response.Bytes())
}

func (u *User) GetAllUsers(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	usersDomain, err := u.UserService.GetAllUsers(contextControl)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, ErrorToGetUsers, err)
		return
	}

	var getUsersResponse []GetUserResponse
	copier.Copy(&getUsersResponse, &usersDomain)
	response := utils.ObjectResponse(getUsersResponse, SuccessToGetUsers)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (u *User) GetUserByID(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	userIDParam := chi.URLParam(r, "id")

	if userIDParam == "" {
		utils.HandleError(w, u.LoggerSugar, http.StatusBadRequest, MissingUserIDParameter, errors.New(UserIDIsRequired))
		return
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusBadRequest, ErrorToParseUser, err)
		return
	}

	userDomain, err := u.UserService.GetUserByID(contextControl, userID)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, ErrorToGetUser, err)
		return
	}

	var getUserResponse []GetUserResponse
	copier.Copy(&getUserResponse, &userDomain)
	response := utils.ObjectResponse(getUserResponse, SuccessToGetUser)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (u *User) UpdateUser(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	userIDParam := chi.URLParam(r, "id")

	if userIDParam == "" {
		utils.HandleError(w, u.LoggerSugar, http.StatusBadRequest, MissingUserIDParameter, errors.New(UserIDIsRequired))
		return
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusBadRequest, ErrorToParseUser, err)
		return
	}

	userIDToken, err := middlewares.ExtractUserID(r)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusUnauthorized, ErrorToExtractUserID, err)
		return
	}

	if userID != userIDToken {
		utils.HandleError(w, u.LoggerSugar, http.StatusForbidden, ErrUserPermissionDenied, errors.New(ErrUserPermissionDenied))
		return
	}

	var userRequest UserRequest

	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, ErrorToDecodeUserRequest, err)
		return
	}

	userRequest.ID = userID

	if err := userRequest.prepareUser("edit"); err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, ErrorToPrepareUser, err)
		return
	}

	var userDomain domain.UserDomain
	copier.Copy(&userDomain, &userRequest)

	userDomain, err = u.UserService.UpdateUser(contextControl, userDomain)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, ErrorToUpdateUser, err)
		return
	}

	var userResponse UserResponse
	copier.Copy(&userResponse, &userDomain)

	response := utils.ObjectResponse(userResponse, SuccessToUpdateUser)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (u *User) SoftDeleteUser(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	userIDParam := chi.URLParam(r, "id")

	if userIDParam == "" {
		utils.HandleError(w, u.LoggerSugar, http.StatusBadRequest, MissingUserIDParameter, errors.New(UserIDIsRequired))
		return
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusBadRequest, ErrorToParseUser, err)
		return
	}

	userIDToken, err := middlewares.ExtractUserID(r)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusUnauthorized, ErrorToExtractUserID, err)
		return
	}

	if userID != userIDToken {
		utils.HandleError(w, u.LoggerSugar, http.StatusForbidden, ErrUserPermissionDenied, errors.New(ErrUserPermissionDenied))
		return
	}

	err = u.UserService.SoftDeleteUser(contextControl, userID)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, ErrorToUpdateUser, err)
		return
	}

	response := utils.ObjectResponse(nil, SuccessToDeleteUser)
	utils.ResponseReturn(w, http.StatusNoContent, response.Bytes())
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var loginUserRequest LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginUserRequest); err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, ErrorToDecodeLoginRequest, err)
		return
	}

	var userDomain domain.UserDomain
	copier.Copy(&userDomain, &loginUserRequest)

	userDB, err := u.UserService.GetUserByUsername(contextControl, userDomain)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, ErrorToGetUserByUsername, err)
		return
	}

	if err = middlewares.VerifyPassword(userDB.Password, userDomain.Password); err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusUnauthorized, ErrorToVerifyPassword, err)
		return
	}

	token, err := security.CreateToken(userDB.ID)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, ErrorToCreateToken, err)
		return
	}

	//If you want to use the function: auth.extractTokenFromBearer
	// you don't need to set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true, // It's not possible to access the cookie via JavaScript
		Secure:   true, // Only send the cookie if the request is being sent over HTTPS
		SameSite: http.SameSiteStrictMode,
	})

	var loginUserResponse LoginUserResponse
	copier.Copy(&loginUserResponse, &userDomain)

	loginUserResponse.Token = token

	response := utils.ObjectResponse(loginUserResponse, SuccessToLogin)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (userRequest *UserRequest) prepareUser(step string) error {
	if err := userRequest.validateUser(step); err != nil {
		return err
	}

	if err := userRequest.formatUser(step); err != nil {
		return err
	}

	return nil
}

func (userRequest *UserRequest) validateUser(step string) error {
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

func (userRequest *UserRequest) formatUser(step string) error {
	userRequest.Name = strings.TrimSpace(userRequest.Name)
	userRequest.Username = strings.TrimSpace(userRequest.Username)
	userRequest.Email = strings.TrimSpace(userRequest.Email)

	if step == "register" {
		hashedPassword, err := middlewares.Hash(userRequest.Password)
		if err != nil {
			return err
		}

		userRequest.Password = string(hashedPassword)
	}
	return nil
}
