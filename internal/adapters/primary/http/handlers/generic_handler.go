package handlers

import (
	msg "github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/messages"
	"net/http"

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
