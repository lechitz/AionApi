package handlers

import (
	constants "github.com/lechitz/AionApi/internal/adapters/primary/http/constants"
	"net/http"

	"github.com/lechitz/AionApi/pkg/utils"
	"go.uber.org/zap"
)

type Generic struct {
	LoggerSugar *zap.SugaredLogger
}

func (h *Generic) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	utils.ResponseReturn(w, http.StatusOK, []byte(constants.MsgServiceIsHealthy))
	h.LoggerSugar.Infow(constants.MsgServiceIsHealthy)
}

func (h *Generic) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	utils.ResponseReturn(w, http.StatusNotFound, []byte(constants.MsgResourceNotFound))
	h.LoggerSugar.Infow(constants.MsgResourceNotFound, "path", r.URL.Path)
}
