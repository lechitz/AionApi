// Package user handlers provides HTTP handlers for user-related endpoints.
package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user/dto"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user/constants"
	"github.com/lechitz/AionApi/internal/core/domain" //TODO: ajustar !!!!
)

// parseUserIDParam extracts and validates the user ID parameter from the URL.
func parseUserIDParam(r *http.Request, log output.ContextLogger) (uint64, error) {
	userIDParam := chi.URLParam(r, commonkeys.UserID)

	if userIDParam == "" {
		err := sharederrors.NewValidationError(commonkeys.UserID, constants.ErrMissingUserIDParam)
		log.Errorw(constants.ErrMissingUserIDParam, commonkeys.Error, err.Error())
		return 0, err
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		validationErr := sharederrors.NewValidationError(commonkeys.UserID, constants.ErrInvalidUserIDParam)
		log.Errorw(constants.ErrInvalidUserIDParam, commonkeys.Error, validationErr.Error())
		return 0, validationErr
	}

	return userID, nil
}

// buildUserDomainFromUpdate builds a user domain from an update user request.
func buildUserDomainFromUpdate(userID uint64, req dto.UpdateUserRequest) domain.User {
	userDomain := domain.User{ID: userID}
	if req.Name != nil {
		userDomain.Name = *req.Name
	}
	if req.Username != nil {
		userDomain.Username = *req.Username
	}
	if req.Email != nil {
		userDomain.Email = *req.Email
	}
	return userDomain
}

// parseUpdateUserRequest parses an update user request.
func (h *Handler) parseUpdateUserRequest(r *http.Request) (dto.UpdateUserRequest, error) {
	var req dto.UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// writeUpdateSuccess writes a success response for a user update operation.
func (h *Handler) writeUpdateSuccess(w http.ResponseWriter, span trace.Span, userUpdated domain.User) {
	span.SetAttributes(attribute.String(commonkeys.UpdatedUsername, userUpdated.Username))
	span.SetStatus(codes.Ok, constants.StatusUserUpdated)
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK))

	res := dto.UpdateUserResponse{
		ID:        userUpdated.ID,
		Name:      &userUpdated.Name,
		Username:  &userUpdated.Username,
		Email:     &userUpdated.Email,
		UpdatedAt: userUpdated.UpdatedAt,
	}

	httpresponse.WriteSuccess(w, http.StatusOK, res, constants.MsgUserUpdated)
}
