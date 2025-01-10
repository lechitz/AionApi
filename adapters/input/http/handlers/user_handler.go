package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/adapters/input/constants"
	"github.com/lechitz/AionApi/adapters/input/http/dto"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/pkg/utils"
	"github.com/lechitz/AionApi/ports/input"
	"go.uber.org/zap"
	"net/http"
)

type User struct {
	UserService input.IUserService
	LoggerSugar *zap.SugaredLogger
}

func (u *User) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	var createUserRequest dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusBadRequest, constants.ErrorToDecodeUserRequest, err)
		return
	}

	var userDomain domain.UserDomain
	copier.Copy(&userDomain, &createUserRequest)

	user, err := u.UserService.CreateUser(contextControl, userDomain)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, constants.ErrorToCreateUser, err)
		return
	}

	var createUserResponse dto.CreateUserResponse
	copier.Copy(&createUserResponse, &user)

	response := utils.ObjectResponse(createUserResponse, constants.SuccessToCreateUser)
	utils.ResponseReturn(w, http.StatusCreated, response.Bytes())
}

func (u *User) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	usersDomain, err := u.UserService.GetAllUsers(contextControl)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, constants.ErrorToGetUsers, err)
		return
	}

	var getUsersResponse []dto.GetUserResponse
	copier.Copy(&getUsersResponse, &usersDomain)
	response := utils.ObjectResponse(getUsersResponse, constants.SuccessToGetUsers)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (u *User) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	userID, err := utils.UserIDFromParam(w, u.LoggerSugar, r)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusBadRequest, constants.ErrorToParseUser, err)
		return
	}

	userDomain, err := u.UserService.GetUserByID(contextControl, userID)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, constants.ErrorToGetUser, err)
		return
	}

	var getUserResponse []dto.GetUserResponse
	copier.Copy(&getUserResponse, &userDomain)
	response := utils.ObjectResponse(getUserResponse, constants.SuccessToGetUser)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (u *User) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	userIDParam, err := utils.UserIDFromParam(w, u.LoggerSugar, r)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusBadRequest, constants.ErrorToParseUser, err)
		return
	}

	userIDToken, err := getUserIDFromContext(r.Context())
	if err != nil {
		u.LoggerSugar.Errorw("Failed to extract userID from context", "error", err.Error())
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if userIDParam != userIDToken {
		utils.HandleError(w, u.LoggerSugar, http.StatusForbidden, constants.ErrorUserPermissionDenied, errors.New(constants.ErrorUserPermissionDenied))
		return
	}

	var updateUserRequest dto.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&updateUserRequest); err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, constants.ErrorToDecodeUserRequest, err)
		return
	}

	updateUserRequest.ID = userIDParam

	var userDomain domain.UserDomain
	copier.Copy(&userDomain, &updateUserRequest)

	userDomain, err = u.UserService.UpdateUser(contextControl, userDomain)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, constants.ErrorToUpdateUser, err)
		return
	}

	var updateUserResponse dto.UpdateUserResponse
	copier.Copy(&updateUserResponse, &userDomain)

	response := utils.ObjectResponse(updateUserResponse, constants.SuccessToUpdateUser)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (u *User) SoftDeleteUserHandler(w http.ResponseWriter, r *http.Request) {

	contextControl := domain.ContextControl{
		Context: context.Background(),
	}

	userIDParam, err := utils.UserIDFromParam(w, u.LoggerSugar, r)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusBadRequest, constants.ErrorToParseUser, err)
		return
	}

	userIDToken, err := getUserIDFromContext(r.Context())
	if err != nil {
		u.LoggerSugar.Errorw("Failed to extract userID from context", "error", err.Error())
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if userIDParam != userIDToken {
		utils.HandleError(w, u.LoggerSugar, http.StatusForbidden, constants.ErrorUserPermissionDenied, errors.New(constants.ErrorUserPermissionDenied))
		return
	}

	err = u.UserService.SoftDeleteUser(contextControl, userIDParam)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, constants.ErrorToUpdateUser, err)
		return
	}

	response := utils.ObjectResponse(nil, constants.SuccessToDeleteUser)
	utils.ResponseReturn(w, http.StatusNoContent, response.Bytes())
}

func getUserIDFromContext(ctx context.Context) (uint64, error) {
	userID, ok := ctx.Value("id").(uint64)
	if !ok {
		return 0, errors.New(constants.ErrorToExtractUserIDFromContext)
	}
	return userID, nil
}
