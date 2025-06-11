// Package validator provides utility functions for validating user parameters.
package validator

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

// missingUserIDParam indicates the error message when the user_id URL parameter is absent.
const missingUserIDParam = "missing user ID parameter"

// errorParsingUserID indicates the error message when parsing user_id fails.
const errorParsingUserID = "error parsing user ID"

// userIDRequired indicates the error message when user_id must be provided.
const userIDRequired = "user ID is required"

// UserID is the name of the URL parameter key for user identification.
const UserID = "user_id"

// ParseUserIDParam extracts and validates the user ID parameter from the URL, parses it into uint64, and handles any parsing errors.
func ParseUserIDParam(w http.ResponseWriter, r *http.Request, log logger.Logger) (uint64, error) {
	userIDParam := chi.URLParam(r, UserID)

	if userIDParam == "" {
		err := errors.New(userIDRequired)
		response.HandleError(w, log, http.StatusBadRequest, missingUserIDParam, err)
		return 0, err
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		response.HandleError(w, log, http.StatusBadRequest, errorParsingUserID, err)
		return 0, err
	}

	return userID, nil
}
