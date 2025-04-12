package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

func ResponseReturn(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if len(body) != 0 {
		w.Write(body)
	}
}

func ObjectResponse(obj any, message string) *bytes.Buffer {
	response := struct {
		Message string    `json:"message,omitempty"`
		Result  any       `json:"result,omitempty"`
		Date    time.Time `json:"date,omitempty"`
	}{
		Message: message,
		Result:  obj,
		Date:    time.Now().UTC(),
	}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(response)
	return body
}

func HandleError(w http.ResponseWriter, logger logger.Logger, status int, msg string, err error) {
	if err != nil {
		logger.Errorw("operation failed",
			"message", msg,
			"error", err.Error(),
			"status", status,
		)
		response := ObjectResponse(nil, msg+": "+err.Error())
		ResponseReturn(w, status, response.Bytes())
	} else {
		logger.Warnw("operation returned warning",
			"message", msg,
			"status", status,
		)
		response := ObjectResponse(nil, msg)
		ResponseReturn(w, status, response.Bytes())
	}
}

func HandleCriticalError(logger logger.Logger, message string, err error) {
	if err != nil {
		logger.Errorw("critical failure",
			"message", message,
			"error", err.Error(),
		)
		panic(err)
	} else {
		logger.Errorw("critical failure",
			"message", message,
		)
		panic(message)
	}
}
