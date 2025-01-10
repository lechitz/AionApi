package utils

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/internal/core/constants"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	ErrorParsingUserId     = "error parsing user ID"
	MissingUserIdParameter = "missing user ID parameter"
)

func UserIDFromParam(w http.ResponseWriter, logger *zap.SugaredLogger, r *http.Request) (uint64, error) {
	userIDParam := chi.URLParam(r, "id")

	if userIDParam == "" {
		HandleError(w, logger, http.StatusBadRequest, MissingUserIdParameter, errors.New(constants.UserIdIsRequired))
		return 0, errors.New(constants.UserIdIsRequired)
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		HandleError(w, logger, http.StatusBadRequest, ErrorParsingUserId, errors.New(constants.ErrorParsingUserId))
		return 0, errors.New(ErrorParsingUserId)
	}

	return userID, nil
}
