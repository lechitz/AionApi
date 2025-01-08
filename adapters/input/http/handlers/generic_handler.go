package handlers

import (
	"github.com/lechitz/AionApi/adapters/input/constants"
	"github.com/lechitz/AionApi/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

type Generic struct {
	LoggerSugar *zap.SugaredLogger
}

func (h *Generic) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	utils.HandleError(w, h.LoggerSugar, http.StatusInternalServerError, constants.ErrorToHealthCheck, nil)
	utils.ResponseReturn(w, http.StatusOK, nil)
}

func (h *Generic) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	utils.HandleError(w, h.LoggerSugar, http.StatusNotFound, constants.ErrorNotFoundHandler, nil)
	utils.ResponseReturn(w, http.StatusNotFound, nil)
}
