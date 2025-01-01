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

func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var userRequest UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		u.LoggerSugar.Errorw(utils.ErrorToDecodeUserRequest, "error", err.Error())
		response := utils.ObjectResponse(utils.ErrorToDecodeUserRequest, err.Error())
		utils.ResponseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	if err := userRequest.prepareUser("register"); err != nil {
		u.LoggerSugar.Errorw(utils.ErrorToPrepareUser, "error", err.Error())
		response := utils.ObjectResponse(utils.ErrorToPrepareUser, err.Error())
		utils.ResponseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var userDomain domain.UserDomain
	copier.Copy(&userDomain, &userRequest)

	userDomain, err := u.UserService.CreateUser(contextControl, userDomain)
	if err != nil {
		u.LoggerSugar.Errorw(utils.ErrorToCreateUser, "error", err.Error())
		response := utils.ObjectResponse(utils.ErrorToCreateUser, err.Error())
		utils.ResponseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var userResponse UserResponse
	copier.Copy(&userResponse, &userDomain)
	response := utils.ObjectResponse(userResponse, utils.SuccessToCreateUser)
	utils.ResponseReturn(w, http.StatusCreated, response.Bytes())
}

func (u *User) GetAllUsers(w http.ResponseWriter) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	usersDomain, err := u.UserService.GetAllUsers(contextControl)
	if err != nil {
		u.LoggerSugar.Errorw(utils.ErrorToGetUsers, "error", err.Error())
		response := utils.ObjectResponse(utils.ErrorToGetUsers, err.Error())
		utils.ResponseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var getUsersResponse []GetUserResponse
	copier.Copy(&getUsersResponse, &usersDomain)
	response := utils.ObjectResponse(getUsersResponse, utils.SucessToGetUsers)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (u *User) GetUserByID(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	userIDParam := chi.URLParam(r, "id")

	if userIDParam == "" {
		u.LoggerSugar.Errorw(utils.MissingUserIDParameter)
		response := utils.ObjectResponse(utils.MissingUserIDParameter, utils.UserIDIsRequired)
		utils.ResponseReturn(w, http.StatusBadRequest, response.Bytes())
		return
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		u.LoggerSugar.Errorw(utils.ErrorToParseUser, "error", err.Error())
		response := utils.ObjectResponse(utils.ErrorToParseUser, err.Error())
		utils.ResponseReturn(w, http.StatusBadRequest, response.Bytes())
		return
	}

	userDomain, err := u.UserService.GetUserByID(contextControl, userID)
	if err != nil {
		u.LoggerSugar.Errorw(utils.ErrorToGetUser, "error", err.Error())
		response := utils.ObjectResponse(utils.ErrorToGetUser, err.Error())
		utils.ResponseReturn(w, http.StatusInternalServerError, response.Bytes())
		return
	}

	var getUserResponse []GetUserResponse
	copier.Copy(&getUserResponse, &userDomain)
	response := utils.ObjectResponse(getUserResponse, utils.SucessToGetUser)
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
		return errors.New(utils.NameIsRequired)
	}
	if userRequest.Username == "" {
		return errors.New(utils.UsernameIsRequired)
	}
	if userRequest.Email == "" {
		return errors.New(utils.EmailIsRequired)
	}

	if err := checkmail.ValidateFormat(userRequest.Email); err != nil {
		return errors.New(utils.InvalidEmail)
	}

	if step == "register" && userRequest.Password == "" {
		return errors.New(utils.PasswordIsRequired)
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
