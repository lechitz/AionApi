package handlers

import (
	"net/http"

	"github.com/lechitz/AionApi/adapters/primary/http/constants"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/response"

	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

type Generic struct {
	Logger logger.Logger
}

func NewGeneric(logger logger.Logger) *Generic {
	return &Generic{Logger: logger}
}

func (h *Generic) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response.ResponseReturn(w, http.StatusOK, []byte(constants.MsgServiceIsHealthy))
}

func (h *Generic) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	response.ResponseReturn(w, http.StatusNotFound, []byte(constants.MsgResourceNotFound))
	h.Logger.Infow(constants.MsgResourceNotFound, "path", r.URL.Path)
}
