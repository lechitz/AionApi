package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/otel"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/dto"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/utils/validator"

	"github.com/lechitz/AionApi/internal/core/domain"
	inputHttp "github.com/lechitz/AionApi/internal/core/ports/input/http"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"

	"github.com/jinzhu/copier"
)

// User represents a handler for managing user-related operations and dependencies.
// It combines user service functionality and logging capabilities.
type User struct {
	UserService inputHttp.UserService
	Logger      logger.Logger
}

// NewUser initializes and returns a new User instance with provided user service and logger dependencies.
func NewUser(userService inputHttp.UserService, logger logger.Logger) *User {
	return &User{
		UserService: userService,
		Logger:      logger,
	}
}

// CreateUserHandler handles HTTP POST requests to create a new user based on the provided request payload.
// It returns a response with the newly created user's data on success or an appropriate error response on failure.
// The request payload must contain the following fields: username, email, password.
// The password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character.
// The username must be at least 3 characters long and contain only alphanumeric characters.
// The email must be a valid email address.
// The username and email must be unique.
// The username and email must not be empty.
func (u *User) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.TracerCreateUserHandler)
	defer span.End()

	span.AddEvent("decoding request")
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusBadRequest, constants.ErrorToDecodeUserRequest, err)
		return
	}
	span.SetAttributes(
		attribute.String(constants.Username, req.Username),
		attribute.String(constants.Email, req.Email),
	)

	var userDomain domain.UserDomain
	_ = copier.Copy(&userDomain, &req)

	span.AddEvent("calling UserService.CreateUser")
	user, err := u.UserService.CreateUser(ctx, userDomain, req.Password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToCreateUser, err)
		return
	}
	span.SetAttributes(attribute.String(constants.UserID, strconv.FormatUint(user.ID, 10)))
	span.SetStatus(codes.Ok, "User created")

	var res dto.CreateUserResponse
	_ = copier.Copy(&res, &user)

	body := response.ObjectResponse(res, constants.SuccessToCreateUser, u.Logger)
	response.Return(w, http.StatusCreated, body.Bytes(), u.Logger)
}

// GetAllUsersHandler handles HTTP requests to retrieve all users and returns the data in the response.// GetAllUsersHandler handles HTTP GET requests to retrieve all users and returns the data as a response.
func (u *User) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.TracerGetAllUsersHandler)
	defer span.End()

	span.AddEvent("calling UserService.GetAllUsers")
	users, err := u.UserService.GetAllUsers(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToGetUsers, err)
		return
	}
	span.SetAttributes(attribute.Int("users_count", len(users)))
	span.SetStatus(codes.Ok, "Users retrieved")

	var res []dto.GetUserResponse
	_ = copier.Copy(&res, &users)

	body := response.ObjectResponse(res, constants.SuccessToGetUsers, u.Logger)
	response.Return(w, http.StatusOK, body.Bytes(), u.Logger)
}

// GetUserByIDHandler handles HTTP requests to retrieve a user by their ID and returns the user's data in the response.
func (u *User) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerCreateUserHandler).
		Start(r.Context(), constants.TracerGetUserByIDHandler)
	defer span.End()

	userID, err := validator.ParseUserIDParam(w, r, u.Logger)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusBadRequest, constants.ErrorToParseUser, err)
		return
	}
	span.SetAttributes(attribute.String(constants.UserID, strconv.FormatUint(userID, 10)))

	user, err := u.UserService.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToGetUser, err)
		return
	}
	span.SetAttributes(attribute.String(constants.Username, user.Username))
	span.SetStatus(codes.Ok, "User retrieved")

	res := dto.GetUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	body := response.ObjectResponse(res, constants.SuccessToGetUser, u.Logger)
	response.Return(w, http.StatusOK, body.Bytes(), u.Logger)
}

// UpdateUserHandler handles HTTP PUT requests to update an existing user's data based on the provided request payload.
func (u *User) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.TracerUpdateUserHandler)
	defer span.End()

	userID, ok := ctx.Value(constants.UserID).(uint64)
	if !ok {
		span.RecordError(errMissingUserID())
		span.SetStatus(codes.Error, "missing user id in context")
		u.logAndHandleError(
			w,
			http.StatusUnauthorized,
			constants.ErrorUnauthorizedAccessMissingToken,
			nil,
		)
		return
	}
	span.SetAttributes(attribute.String(constants.UserID, strconv.FormatUint(userID, 10)))

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
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

	span.AddEvent("calling UserService.UpdateUser")
	userUpdated, err := u.UserService.UpdateUser(ctx, userDomain)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToUpdateUser, err)
		return
	}
	span.SetAttributes(attribute.String("updated_username", userUpdated.Username))
	span.SetStatus(codes.Ok, "User updated")

	res := dto.UpdateUserResponse{
		ID:       userUpdated.ID,
		Name:     &userUpdated.Name,
		Username: &userUpdated.Username,
		Email:    &userUpdated.Email,
	}

	body := response.ObjectResponse(res, constants.SuccessToUpdateUser, u.Logger)
	response.Return(w, http.StatusOK, body.Bytes(), u.Logger)
}

// UpdatePasswordHandler handles the HTTP request to update a user's password and refreshes their authentication token.
func (u *User) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.TracerUpdatePasswordHandler)
	defer span.End()

	var req dto.UpdatePasswordUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusBadRequest, constants.ErrorToDecodeUserRequest, err)
		return
	}

	userID, ok := ctx.Value(constants.UserID).(uint64)
	if !ok {
		span.RecordError(errMissingUserID())
		span.SetStatus(codes.Error, "missing user id in context")
		u.logAndHandleError(
			w,
			http.StatusUnauthorized,
			constants.ErrorUnauthorizedAccessMissingToken,
			nil,
		)
		return
	}
	span.SetAttributes(attribute.String(constants.UserID, strconv.FormatUint(userID, 10)))

	clearAuthCookie(w)

	userDomain := domain.UserDomain{ID: userID}
	span.AddEvent("calling UserService.UpdateUserPassword")
	_, newToken, err := u.UserService.UpdateUserPassword(ctx, userDomain, req.Password, req.NewPassword)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToUpdateUser, err)
		return
	}
	setAuthCookie(w, newToken, 0)
	span.SetStatus(codes.Ok, "Password updated")

	body := response.ObjectResponse(nil, constants.SuccessToUpdatePassword, u.Logger)
	response.Return(w, http.StatusOK, body.Bytes(), u.Logger)
}

// SoftDeleteUserHandler handles the soft deletion of a user by ID extracted from the request context.
// Responds with HTTP 204 on success or appropriate error response if the operation fails.
func (u *User) SoftDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.TracerSoftDeleteUserHandler)
	defer span.End()

	userID, ok := ctx.Value(constants.UserID).(uint64)
	if !ok {
		span.RecordError(errMissingUserID())
		span.SetStatus(codes.Error, "missing user id in context")
		u.logAndHandleError(
			w,
			http.StatusUnauthorized,
			constants.ErrorUnauthorizedAccessMissingToken,
			nil,
		)
		return
	}
	span.SetAttributes(attribute.String(constants.UserID, strconv.FormatUint(userID, 10)))

	span.AddEvent("calling UserService.SoftDeleteUser")
	if err := u.UserService.SoftDeleteUser(ctx, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToSoftDeleteUser, err)
		return
	}
	clearAuthCookie(w)
	span.SetStatus(codes.Ok, "User soft deleted")

	body := response.ObjectResponse(nil, constants.SuccessUserSoftDeleted, u.Logger)
	response.Return(w, http.StatusNoContent, body.Bytes(), u.Logger)
}

// errMissingUserID returns an error indicating that the user ID is missing from the request context.
func errMissingUserID() error {
	return &MissingUserIDError{}
}

// MissingUserIDError is an error type indicating that the user ID is missing from the request context.
type MissingUserIDError struct{}

func (e *MissingUserIDError) Error() string {
	return "userID missing from context"
}

// logAndHandleError logs the error with a message and sends an HTTP error response to the client.
func (u *User) logAndHandleError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		u.Logger.Errorw(message, constants.Error, err.Error())
	} else {
		u.Logger.Errorw(message)
	}
	response.HandleError(w, u.Logger, status, message, err)
}
