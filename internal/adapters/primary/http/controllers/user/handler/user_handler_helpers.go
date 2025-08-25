// Package handler (user) controllers provide HTTP controllers for user-related endpoints.
package handler

import (
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

	"github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/user/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/user/dto"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
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

// writeUpdateSuccess writes a success response for a user update operation using a DTO.
// This avoids coupling controllers to domain entities.
func (h *Handler) writeUpdateSuccess(w http.ResponseWriter, span trace.Span, res dto.UpdateUserResponse) {
	updatedUsername := ""
	if res.Username != nil {
		updatedUsername = *res.Username
	}

	span.SetAttributes(
		attribute.String(commonkeys.UpdatedUsername, updatedUsername),
		attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusOK),
	)
	span.SetStatus(codes.Ok, constants.StatusUserUpdated)

	httpresponse.WriteSuccess(w, http.StatusOK, res, constants.MsgUserUpdated)
}
