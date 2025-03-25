package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/adapters/input/http/dto"
	msg "github.com/lechitz/AionApi/adapters/input/http/handlers/messages"
	"github.com/lechitz/AionApi/core/domain"
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
	ctx := domain.ContextControl{BaseContext: r.Context()}

	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logAndHandleError(w, http.StatusBadRequest, msg.ErrorToDecodeUserRequest, err)
		return
	}

	var userDomain domain.UserDomain
	_ = copier.Copy(&userDomain, &req)

	user, err := u.UserService.CreateUser(ctx, userDomain, req.Password)
	if err != nil {
		u.logAndHandleError(w, http.StatusInternalServerError, msg.ErrorToCreateUser, err)
		return
	}

	var res dto.CreateUserResponse
	_ = copier.Copy(&res, &user)

	u.LoggerSugar.Infow(msg.SuccessToCreateUser, contextkeys.Username, res.Username)
	utils.ResponseReturn(w, http.StatusCreated, utils.ObjectResponse(res, msg.SuccessToCreateUser).Bytes())
}

func (u *User) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := domain.ContextControl{BaseContext: r.Context()}

	users, err := u.UserService.GetAllUsers(ctx)
	if err != nil {
		u.logAndHandleError(w, http.StatusInternalServerError, msg.ErrorToGetUsers, err)
		return
	}

	var res []dto.GetUserResponse
	_ = copier.Copy(&res, &users)

	u.LoggerSugar.Infow(msg.SuccessToGetUsers, contextkeys.Users, res)
	utils.ResponseReturn(w, http.StatusOK, utils.ObjectResponse(res, msg.SuccessToGetUsers).Bytes())
}

func (u *User) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := domain.ContextControl{BaseContext: r.Context()}

	userID, err := utils.UserIDFromParam(w, u.LoggerSugar, r)
	if err != nil {
		u.logAndHandleError(w, http.StatusBadRequest, msg.ErrorToParseUser, err)
		return
	}

	user, err := u.UserService.GetUserByID(ctx, userID)
	if err != nil {
		u.logAndHandleError(w, http.StatusInternalServerError, msg.ErrorToGetUser, err)
		return
	}

	res := dto.GetUserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
	}

	u.LoggerSugar.Infow(msg.SuccessToGetUser, contextkeys.User, res)
	utils.ResponseReturn(w, http.StatusOK, utils.ObjectResponse(res, msg.SuccessToGetUser).Bytes())
}

func (u *User) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := domain.ContextControl{BaseContext: r.Context()}

	userID, ok := ctx.BaseContext.Value(contextkeys.UserID).(uint64)
	if !ok {
		u.logAndHandleError(w, http.StatusUnauthorized, msg.ErrorUnauthorizedAccessMissingToken, nil)
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logAndHandleError(w, http.StatusBadRequest, msg.ErrorToDecodeUserRequest, err)
		return
	}

	userDomain := domain.UserDomain{ID: userID}
	if req.Name != nil {
		userDomain.Name = *req.Name
	}
	if req.Username != nil {
		userDomain.Username = *req.Username
	}
	if req.Email != nil {
		userDomain.Email = *req.Email
	}

	userUpdated, err := u.UserService.UpdateUser(ctx, userDomain)
	if err != nil {
		u.logAndHandleError(w, http.StatusInternalServerError, msg.ErrorToUpdateUser, err)
		return
	}

	res := dto.UpdateUserResponse{
		ID:       userUpdated.ID,
		Name:     &userUpdated.Name,
		Username: &userUpdated.Username,
		Email:    &userUpdated.Email,
	}

	u.LoggerSugar.Infow(msg.SuccessToUpdateUser, contextkeys.Username, res.Username)
	utils.ResponseReturn(w, http.StatusOK, utils.ObjectResponse(res, msg.SuccessToUpdateUser).Bytes())
}

func (u *User) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := domain.ContextControl{BaseContext: r.Context()}

	var req dto.UpdatePasswordUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logAndHandleError(w, http.StatusBadRequest, msg.ErrorToDecodeUserRequest, err)
		return
	}

	userID, ok := ctx.BaseContext.Value(contextkeys.UserID).(uint64)
	if !ok {
		u.logAndHandleError(w, http.StatusUnauthorized, msg.ErrorUnauthorizedAccessMissingToken, nil)
		return
	}

	clearAuthCookie(w)

	userDomain := domain.UserDomain{ID: userID}
	_, newToken, err := u.UserService.UpdateUserPassword(ctx, userDomain, req.Password, req.NewPassword)
	if err != nil {
		u.logAndHandleError(w, http.StatusInternalServerError, msg.ErrorToUpdateUser, err)
		return
	}

	setAuthCookie(w, newToken, 0)

	u.LoggerSugar.Infow(msg.SuccessToUpdatePassword, contextkeys.UserID, userID)
	utils.ResponseReturn(w, http.StatusOK, utils.ObjectResponse(nil, msg.SuccessToUpdatePassword).Bytes())
}

func (u *User) SoftDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := domain.ContextControl{BaseContext: r.Context()}

	userID, ok := ctx.BaseContext.Value(contextkeys.UserID).(uint64)
	if !ok {
		u.logAndHandleError(w, http.StatusUnauthorized, msg.ErrorUnauthorizedAccessMissingToken, nil)
		return
	}

	if err := u.UserService.SoftDeleteUser(ctx, userID); err != nil {
		u.logAndHandleError(w, http.StatusInternalServerError, msg.ErrorToSoftDeleteUser, err)
		return
	}

	clearAuthCookie(w)

	u.LoggerSugar.Infow(msg.SuccessUserSoftDeleted, contextkeys.UserID, userID)
	utils.ResponseReturn(w, http.StatusNoContent, utils.ObjectResponse(nil, msg.SuccessUserSoftDeleted).Bytes())
}

func (u *User) logAndHandleError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		u.LoggerSugar.Errorw(message, contextkeys.Error, err.Error())
	} else {
		u.LoggerSugar.Errorw(message)
	}
	utils.HandleError(w, u.LoggerSugar, status, message, err)
}
