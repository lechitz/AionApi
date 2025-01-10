package handlers

import (
	"github.com/lechitz/AionApi/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

type Generic struct {
	LoggerSugar *zap.SugaredLogger
}

func (h *Generic) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	utils.ResponseReturn(w, http.StatusOK, []byte(`{"message": "Service is healthy"}`))
}

func (h *Generic) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	utils.ResponseReturn(w, http.StatusNotFound, []byte(`{"error": "Resource not found"}`))
}
