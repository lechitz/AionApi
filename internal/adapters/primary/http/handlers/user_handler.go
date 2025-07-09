package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lechitz/AionApi/internal/platform/config"

	"github.com/lechitz/AionApi/internal/shared/httputils"

	"github.com/lechitz/AionApi/internal/shared/commonkeys"

	"github.com/lechitz/AionApi/internal/shared/ctxkeys"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/otel"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/dto"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/utils/validator"

	"github.com/jinzhu/copier"
)

// TODO: Separar os Handlers e o construtor.. achar um nome bom pro arquivo de NewUser !

// TODO: Melhorar as msgs de Span, Ajustar os erros e logs para ficarem completos !

// User represents a handler for managing user-related operations and dependencies.
type User struct {
	UserService input.UserService
	Logger      output.Logger
	Config      *config.Config
}

// NewUser initializes and returns a new User instance with provided user service and zap dependencies.
func NewUser(userService input.UserService, cfg *config.Config, logger output.Logger) *User {
	return &User{
		UserService: userService,
		Config:      cfg,
		Logger:      logger,
	}
}

// CreateUserHandler handles HTTP POST requests to create a new user based on the provided request payload.
func (u *User) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.TracerCreateUserHandler)
	defer span.End()

	span.AddEvent("decoding request") // TODO: ajustar magic string.
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusBadRequest, constants.ErrorToDecodeUserRequest, err)
		return
	}
	span.SetAttributes(
		attribute.String(commonkeys.Username, req.Username),
		attribute.String(commonkeys.UserEmail, req.Email),
	)

	var userDomain domain.UserDomain
	_ = copier.Copy(&userDomain, &req)

	span.AddEvent("calling UserService.CreateUser") // TODO: ajustar magic string.
	user, err := u.UserService.CreateUser(ctx, userDomain, req.Password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToCreateUser, err)
		return
	}
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(user.ID, 10)))
	span.SetStatus(codes.Ok, "User created") // TODO: ajustar magic string.

	var res dto.CreateUserResponse
	_ = copier.Copy(&res, &user)

	body := response.ObjectResponse(res, constants.SuccessToCreateUser, u.Logger)
	response.Return(w, http.StatusCreated, body.Bytes(), u.Logger)
}

// GetAllUsersHandler handles HTTP requests to retrieve all users and returns the data in the response.
func (u *User) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.TracerGetAllUsersHandler)
	defer span.End()

	span.AddEvent("calling UserService.GetAllUsers") // TODO: ajustar magic string.
	users, err := u.UserService.GetAllUsers(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToGetUsers, err)
		return
	}
	span.SetAttributes(attribute.Int(commonkeys.UsersCount, len(users)))
	span.SetStatus(codes.Ok, "Users retrieved") // TODO: ajustar magic string.

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
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))

	user, err := u.UserService.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToGetUser, err)
		return
	}
	span.SetAttributes(attribute.String(commonkeys.Username, user.Username))
	span.SetStatus(codes.Ok, "User retrieved") // TODO: ajustar magic string.

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

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		span.RecordError(errMissingUserID())
		span.SetStatus(codes.Error, "missing user id in context") // TODO: ajustar magic string.
		u.logAndHandleError(
			w,
			http.StatusUnauthorized,
			constants.ErrorUnauthorizedAccessMissingToken,
			nil,
		)
		return
	}
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))

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

	span.AddEvent("calling UserService.UpdateUser") // TODO: ajustar magic string.
	userUpdated, err := u.UserService.UpdateUser(ctx, userDomain)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToUpdateUser, err)
		return
	}
	span.SetAttributes(attribute.String("updated_username", userUpdated.Username)) // TODO: ajustar magic string.
	span.SetStatus(codes.Ok, "User updated")                                       // TODO: ajustar magic string.

	res := dto.UpdateUserResponse{
		ID:        userUpdated.ID,
		Name:      &userUpdated.Name,
		Username:  &userUpdated.Username,
		Email:     &userUpdated.Email,
		UpdatedAt: userUpdated.UpdatedAt,
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

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		span.RecordError(errMissingUserID())
		span.SetStatus(codes.Error, "missing user id in context") // TODO: ajustar magic string.
		u.logAndHandleError(
			w,
			http.StatusUnauthorized,
			constants.ErrorUnauthorizedAccessMissingToken,
			nil,
		)
		return
	}
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))

	httputils.ClearAuthCookie(w, u.Config.Cookie)

	userDomain := domain.UserDomain{ID: userID}
	span.AddEvent("calling UserService.UpdateUserPassword") // TODO: ajustar magic string.
	_, newToken, err := u.UserService.UpdateUserPassword(ctx, userDomain, req.Password, req.NewPassword)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToUpdateUser, err)
		return
	}

	httputils.SetAuthCookie(w, newToken, u.Config.Cookie)
	span.SetStatus(codes.Ok, "UserPassword updated") // TODO: ajustar magic string.

	body := response.ObjectResponse(nil, constants.SuccessToUpdatePassword, u.Logger)
	response.Return(w, http.StatusOK, body.Bytes(), u.Logger)
}

// SoftDeleteUserHandler handles the soft deletion of a user by ID extracted from the request context.
func (u *User) SoftDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(constants.TracerUserHandler).
		Start(r.Context(), constants.TracerSoftDeleteUserHandler)
	defer span.End()

	userID, ok := ctx.Value(ctxkeys.UserID).(uint64)
	if !ok || userID == 0 {
		span.RecordError(errMissingUserID())
		span.SetStatus(codes.Error, "missing user id in context") // TODO: ajustar magic string.
		u.logAndHandleError(
			w,
			http.StatusUnauthorized,
			constants.ErrorUnauthorizedAccessMissingToken,
			nil,
		)
		return
	}
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))

	span.AddEvent("calling UserService.SoftDeleteUser") // TODO: ajustar magic string.
	if err := u.UserService.SoftDeleteUser(ctx, userID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		u.logAndHandleError(w, http.StatusInternalServerError, constants.ErrorToSoftDeleteUser, err)
		return
	}
	httputils.ClearAuthCookie(w, u.Config.Cookie)
	span.SetStatus(codes.Ok, "User soft deleted") // TODO: ajustar magic string.

	body := response.ObjectResponse(nil, constants.SuccessUserSoftDeleted, u.Logger)
	response.Return(w, http.StatusNoContent, body.Bytes(), u.Logger)
}

// TODO: entender essa parte abaixo de errorMissing, se deveria estar aqui..

// errMissingUserID returns an error indicating that the user ID is missing from the request context.
func errMissingUserID() error {
	return &MissingUserIDError{}
}

// MissingUserIDError is an error type indicating that the user ID is missing from the request context.
type MissingUserIDError struct{}

func (e *MissingUserIDError) Error() string {
	return "userID missing from context"
} // TODO: ajustar magic string.

// logAndHandleError logs the error with a message and sends an HTTP error response to the client.
func (u *User) logAndHandleError(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		u.Logger.Errorw(message, commonkeys.Error, err.Error())
	} else {
		u.Logger.Errorw(message)
	}
	response.HandleError(w, u.Logger, status, message, err)
}
