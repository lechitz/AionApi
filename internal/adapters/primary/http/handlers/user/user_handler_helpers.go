// Package user handlers provides HTTP handlers for user-related endpoints.
package user

import (
	"encoding/json"
	"net/http"

	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user/dto"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user/constants"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
)

// Helpers for parsing and writing errors/success — these are candidates for future genericization (can be moved to a base handler or shared).
func buildUserDomainFromUpdate(userID uint64, req dto.UpdateUserRequest) domain.UserDomain {
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
	return userDomain
}

// writeUpdateSuccess writes a success response for a user update operation.
func (h *Handler) writeUpdateSuccess(w http.ResponseWriter, span trace.Span, userUpdated domain.UserDomain) {
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

// parseUpdateUserRequest parses an update user request.
func (h *Handler) parseUpdateUserRequest(r *http.Request) (dto.UpdateUserRequest, error) {
	var req dto.UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}
