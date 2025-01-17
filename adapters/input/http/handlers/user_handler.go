package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/adapters/input/http/dto"
	msg "github.com/lechitz/AionApi/adapters/input/http/handlers/messages"
	"github.com/lechitz/AionApi/core/domain/entities"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"github.com/lechitz/AionApi/pkg/utils"
	inputHttp "github.com/lechitz/AionApi/ports/input/http"
	"go.uber.org/zap"
)

type User struct {
	UserService inputHttp.IUserService
	LoggerSugar *zap.SugaredLogger
}

func (u *User) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	ctx := &entities.ContextControl{
		BaseContext:     r.Context(),
		CancelCauseFunc: nil,
	}

	var createUserRequest dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusBadRequest, msg.ErrorToDecodeUserRequest, err)
		u.LoggerSugar.Errorw(msg.ErrorToDecodeUserRequest, contextkeys.Error, err.Error())
		return
	}

	var userDomain entities.UserDomain
	copier.Copy(&userDomain, &createUserRequest)

	user, err := u.UserService.CreateUser(*ctx, userDomain, createUserRequest.Password)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, msg.ErrorToCreateUser, err)
		u.LoggerSugar.Errorw(msg.ErrorToCreateUser, contextkeys.Error, err.Error())
		return
	}

	var createUserResponse dto.CreateUserResponse
	copier.Copy(&createUserResponse, &user)

	response := utils.ObjectResponse(createUserResponse, msg.SuccessToCreateUser)
	u.LoggerSugar.Infow(msg.SuccessToCreateUser, contextkeys.Username, createUserResponse.Username)
	utils.ResponseReturn(w, http.StatusCreated, response.Bytes())
}

func (u *User) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {

	contextControl := entities.ContextControl{
		BaseContext: r.Context(),
	}

	usersDomain, err := u.UserService.GetAllUsers(contextControl)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, msg.ErrorToGetUsers, err)
		u.LoggerSugar.Errorw(msg.ErrorToGetUsers, contextkeys.Error, err.Error())
		return
	}

	var getUsersResponse []dto.GetUserResponse
	copier.Copy(&getUsersResponse, &usersDomain)

	response := utils.ObjectResponse(getUsersResponse, msg.SuccessToGetUsers)
	u.LoggerSugar.Infow(msg.SuccessToGetUsers, contextkeys.Users, getUsersResponse)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (u *User) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {

	contextControl := entities.ContextControl{
		BaseContext: r.Context(),
	}

	userIDParam, err := utils.UserIDFromParam(w, u.LoggerSugar, r)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusBadRequest, msg.ErrorToParseUser, err)
		u.LoggerSugar.Errorw(msg.ErrorToParseUser, contextkeys.Error, err.Error())
		return
	}

	userDomain, err := u.UserService.GetUserByID(contextControl, userIDParam)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, msg.ErrorToGetUser, err)
		u.LoggerSugar.Errorw(msg.ErrorToGetUser, contextkeys.Error, err.Error())
		return
	}

	getUserResponse := dto.GetUserResponse{
		ID:       userDomain.ID,
		Name:     userDomain.Name,
		Username: userDomain.Username,
		Email:    userDomain.Email,
	}

	response := utils.ObjectResponse(getUserResponse, msg.SuccessToGetUser)
	u.LoggerSugar.Infow(msg.SuccessToGetUser, contextkeys.User, getUserResponse)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (u *User) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(contextkeys.UserID).(uint64)
	if !ok {
		utils.HandleError(w, u.LoggerSugar, http.StatusUnauthorized, msg.ErrorUnauthorizedAccessMissingToken, nil)
		u.LoggerSugar.Errorw(msg.ErrorUnauthorizedAccessMissingToken, contextkeys.Context, r.Context())
		return
	}

	var updateUserRequest dto.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&updateUserRequest); err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, msg.ErrorToDecodeUserRequest, err)
		u.LoggerSugar.Errorw(msg.ErrorToDecodeUserRequest, contextkeys.Error, err.Error())
		return
	}

	contextControl := entities.ContextControl{
		BaseContext: r.Context(),
	}

	existingUser, err := u.UserService.GetUserByID(contextControl, userID)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, msg.ErrorToGetUser, err)
		u.LoggerSugar.Errorw(msg.ErrorToGetUser, contextkeys.Error, err.Error())
		return
	}

	if updateUserRequest.Name != nil {
		existingUser.Name = *updateUserRequest.Name
	}
	if updateUserRequest.Email != nil {
		existingUser.Email = *updateUserRequest.Email
	}
	if updateUserRequest.Username != nil {
		existingUser.Username = *updateUserRequest.Username
	}

	userDomain := entities.UserDomain{
		ID:       userID,
		Name:     existingUser.Name,
		Username: existingUser.Username,
		Email:    existingUser.Email,
	}

	updateUser, err := u.UserService.UpdateUser(contextControl, userDomain)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, msg.ErrorToUpdateUser, err)
		u.LoggerSugar.Errorw(msg.ErrorToUpdateUser, contextkeys.Error, err.Error())
		return
	}

	updateUserResponse := dto.UpdateUserResponse{
		ID:       updateUser.ID,
		Username: &updateUser.Username,
		Email:    &updateUser.Email,
	}

	response := utils.ObjectResponse(updateUserResponse, msg.SuccessToUpdateUser)
	u.LoggerSugar.Infow(msg.SuccessToUpdateUser, contextkeys.Username, updateUserResponse.Username)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (u *User) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {

	contextControl := entities.ContextControl{
		BaseContext: r.Context(),
	}

	var updatePasswordRequest dto.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&updatePasswordRequest); err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusBadRequest, msg.ErrorToDecodeUserRequest, err)
		u.LoggerSugar.Errorw(msg.ErrorToDecodeUserRequest, contextkeys.Error, err.Error())
		return
	}

	userID, ok := contextControl.BaseContext.Value(contextkeys.UserID).(uint64)
	if !ok {
		utils.HandleError(w, u.LoggerSugar, http.StatusUnauthorized, msg.ErrorUnauthorizedAccessMissingToken, nil)
		u.LoggerSugar.Errorw(msg.ErrorUnauthorizedAccessMissingToken, contextkeys.Context, contextControl.BaseContext)
		return
	}

	var userDomain entities.UserDomain
	userDomain.ID = userID

	_, token, err := u.UserService.UpdateUserPassword(contextControl, userDomain, updatePasswordRequest.Password, updatePasswordRequest.NewPassword)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, msg.ErrorToUpdateUser, err)
		u.LoggerSugar.Errorw(msg.ErrorToUpdateUser, contextkeys.Error, err.Error())
		return
	}

	setAuthCookie(w, token, 0)

	response := utils.ObjectResponse(nil, msg.SuccessToUpdatePassword)
	u.LoggerSugar.Infow(msg.SuccessToUpdatePassword, contextkeys.Username, userDomain.Username)
	utils.ResponseReturn(w, http.StatusOK, response.Bytes())
}

func (u *User) SoftDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	contextControl := entities.ContextControl{
		BaseContext: r.Context(),
	}

	userID, ok := contextControl.BaseContext.Value(contextkeys.UserID).(uint64)
	if !ok {
		utils.HandleError(w, u.LoggerSugar, http.StatusUnauthorized, msg.ErrorUnauthorizedAccessMissingToken, nil)
		u.LoggerSugar.Errorw(msg.ErrorUnauthorizedAccessMissingToken, contextkeys.Context, contextControl.BaseContext)
		return
	}

	err := u.UserService.SoftDeleteUser(contextControl, userID)
	if err != nil {
		utils.HandleError(w, u.LoggerSugar, http.StatusInternalServerError, msg.ErrorToSoftDeleteUser, err)
		u.LoggerSugar.Errorw(msg.ErrorToSoftDeleteUser, contextkeys.Error, err.Error(), contextkeys.UserID, userID)
		return
	}

	response := utils.ObjectResponse(nil, msg.SuccessUserSoftDeleted)
	u.LoggerSugar.Infow(msg.SuccessUserSoftDeleted, contextkeys.UserID, userID)
	utils.ResponseReturn(w, http.StatusNoContent, response.Bytes())
}
