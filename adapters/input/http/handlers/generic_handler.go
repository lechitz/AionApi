package handlers

import (
	"net/http"

	msg "github.com/lechitz/AionApi/adapters/input/http/handlers/messages"
	"github.com/lechitz/AionApi/pkg/utils"
	"go.uber.org/zap"
)

type Generic struct {
	LoggerSugar *zap.SugaredLogger
}

func (h *Generic) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	utils.ResponseReturn(w, http.StatusOK, []byte(msg.MsgServiceIsHealthy))
	h.LoggerSugar.Infow(msg.MsgServiceIsHealthy)
}

func (h *Generic) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	utils.ResponseReturn(w, http.StatusNotFound, []byte(msg.MsgResourceNotFound))
	h.LoggerSugar.Infow(msg.MsgResourceNotFound, "path", r.URL.Path)
}
