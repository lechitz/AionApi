package validators

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/pkg/utils"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	MissingUserIdParameter = "missing user ID parameter"
	ErrorParsingUserId     = "error parsing user ID"
	UserIdIsRequired       = "user ID is required"
)

func UserIDFromParam(w http.ResponseWriter, logger *zap.SugaredLogger, r *http.Request) (uint64, error) {
	userIDParam := chi.URLParam(r, "id")

	if userIDParam == "" {
		utils.HandleError(w, logger, http.StatusBadRequest, MissingUserIdParameter, errors.New(UserIdIsRequired))
		return 0, errors.New(UserIdIsRequired)
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		utils.HandleError(w, logger, http.StatusBadRequest, ErrorParsingUserId, errors.New(ErrorParsingUserId))
		return 0, errors.New(ErrorParsingUserId)
	}

	return userID, nil
}
