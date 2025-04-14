package handlers

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/constants"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/pkg/utils"
)

type Generic struct {
	Logger logger.Logger
}

func NewGeneric(logger logger.Logger) *Generic {
	return &Generic{Logger: logger}
}

func (h *Generic) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	utils.ResponseReturn(w, http.StatusOK, []byte(constants.MsgServiceIsHealthy))
}

func (h *Generic) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	utils.ResponseReturn(w, http.StatusNotFound, []byte(constants.MsgResourceNotFound))
	h.Logger.Infow(constants.MsgResourceNotFound, "path", r.URL.Path)
}
