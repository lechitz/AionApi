package utils

import (
	"go.uber.org/zap"
	"net/http"
)

func HandleError(w http.ResponseWriter, logger *zap.SugaredLogger, status int, msg string, err error) {
	logger.Errorw(msg, "error", err.Error())
	response := ObjectResponse(msg, err.Error())
	ResponseReturn(w, status, response.Bytes())
}
