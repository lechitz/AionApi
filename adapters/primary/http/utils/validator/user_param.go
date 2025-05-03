package validator

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"net/http"
	"strconv"
)

const (
	missingUserIDParam = "missing user ID parameter"
	errorParsingUserID = "error parsing user ID"
	userIDRequired     = "user ID is required"
	UserID             = "user_id"
)

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
