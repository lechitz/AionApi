package handlers

import (
	"github.com/lechitz/AionApi/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

type Generic struct {
	LoggerSugar *zap.SugaredLogger
}

func (h *Generic) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.LoggerSugar.Warnw("health check")
	utils.ResponseReturn(w, http.StatusOK, nil)
}

func (h *Generic) NotFound(w http.ResponseWriter, r *http.Request) {
	h.LoggerSugar.Warnw("resource not found")
	utils.ResponseReturn(w, http.StatusNotFound, nil)
}
