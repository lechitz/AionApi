package handlers

import (
	"encoding/json"
	"github.com/lechitz/AionApi/adapters/primary/http/constants"
	"github.com/lechitz/AionApi/adapters/primary/http/dto"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/adapters/primary/http/utils/validator"
	"github.com/lechitz/AionApi/internal/core/domain"
	inputHttp "github.com/lechitz/AionApi/internal/core/ports/input/http"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"

	"net/http"

	"github.com/jinzhu/copier"
)

type User struct {
	UserService inputHttp.UserService
	Logger      logger.Logger
}

func NewUser(userService inputHttp.UserService, logger logger.Logger) *User {
	return &User{
		UserService: userService,
		Logger:      logger,
	}
}

func (u *User) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logAndHandleError(w, http.StatusBadRequest, constants.ErrorToDecodeUserRequest, err)
		return
	}

	var userDomain domain.UserDomain
	_ = copier.Copy(&userDomain, &req)

	user, err := u.UserService.CreateUser(ctx, userDomain, req.Password)
	if err != nil {
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToCreateUser, err)
		return
	}

	var res dto.CreateUserResponse
	_ = copier.Copy(&res, &user)

	response.ResponseReturn(w, http.StatusCreated, response.ObjectResponse(res, constants.SuccessToCreateUser).Bytes())
}

func (u *User) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := u.UserService.GetAllUsers(ctx)
	if err != nil {
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToGetUsers, err)
		return
	}

	var res []dto.GetUserResponse
	_ = copier.Copy(&res, &users)

	response.ResponseReturn(w, http.StatusOK, response.ObjectResponse(res, constants.SuccessToGetUsers).Bytes())
}

func (u *User) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := validator.ParseUserIDParam(w, r, u.Logger)
	if err != nil {
		u.logAndHandleError(w, http.StatusBadRequest, constants.ErrorToParseUser, err)
		return
	}

	user, err := u.UserService.GetUserByID(ctx, userID)
	if err != nil {
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToGetUser, err)
		return
	}

	res := dto.GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	response.ResponseReturn(w, http.StatusOK, response.ObjectResponse(res, constants.SuccessToGetUser).Bytes())
}

func (u *User) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(constants.UserID).(uint64)
	if !ok {
		u.logAndHandleError(w, http.StatusUnauthorized, constants.ErrorUnauthorizedAccessMissingToken, nil)
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logAndHandleError(w, http.StatusBadRequest, constants.ErrorToDecodeUserRequest, err)
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
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToUpdateUser, err)
		return
	}

	res := dto.UpdateUserResponse{
		ID:       userUpdated.ID,
		Name:     &userUpdated.Name,
		Username: &userUpdated.Username,
		Email:    &userUpdated.Email,
	}

	response.ResponseReturn(w, http.StatusOK, response.ObjectResponse(res, constants.SuccessToUpdateUser).Bytes())
}

func (u *User) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.UpdatePasswordUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		u.logAndHandleError(w, http.StatusBadRequest, constants.ErrorToDecodeUserRequest, err)
		return
	}

	userID, ok := ctx.Value(constants.UserID).(uint64)
	if !ok {
		u.logAndHandleError(w, http.StatusUnauthorized, constants.ErrorUnauthorizedAccessMissingToken, nil)
		return
	}

	clearAuthCookie(w)

	userDomain := domain.UserDomain{ID: userID}
	_, newToken, err := u.UserService.UpdateUserPassword(ctx, userDomain, req.Password, req.NewPassword)
	if err != nil {
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToUpdateUser, err)
		return
	}

	setAuthCookie(w, newToken, 0)

	response.ResponseReturn(w, http.StatusOK, response.ObjectResponse(nil, constants.SuccessToUpdatePassword).Bytes())
}

func (u *User) SoftDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value(constants.UserID).(uint64)
	if !ok {
		u.logAndHandleError(w, http.StatusUnauthorized, constants.ErrorUnauthorizedAccessMissingToken, nil)
		return
	}

	if err := u.UserService.SoftDeleteUser(ctx, userID); err != nil {
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToSoftDeleteUser, err)
		return
	}

	clearAuthCookie(w)

	response.ResponseReturn(w, http.StatusNoContent, response.ObjectResponse(nil, constants.SuccessUserSoftDeleted).Bytes())
}

func (u *User) logAndHandleError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		u.Logger.Errorw(message, constants.Error, err.Error())
	} else {
		u.Logger.Errorw(message)
	}
	response.HandleError(w, u.Logger, status, message, err)
}
