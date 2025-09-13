// Package handler (user) controllers provide HTTP controllers for user-related endpoints.
package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// parseUserIDParam extracts and validates the user ID parameter from the URL.
func parseUserIDParam(r *http.Request, log logger.ContextLogger) (uint64, error) {
	userIDParam := chi.URLParam(r, commonkeys.UserID)

	if userIDParam == "" {
		err := sharederrors.NewValidationError(commonkeys.UserID, ErrMissingUserIDParam)
		log.Errorw(ErrMissingUserIDParam, commonkeys.Error, err.Error())
		return 0, err
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		validationErr := sharederrors.NewValidationError(commonkeys.UserID, ErrInvalidUserIDParam)
		log.Errorw(ErrInvalidUserIDParam, commonkeys.Error, validationErr.Error())
		return 0, validationErr
	}

	return userID, nil
}
