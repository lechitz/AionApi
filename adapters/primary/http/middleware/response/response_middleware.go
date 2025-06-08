package response

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

func ResponseReturn(w http.ResponseWriter, statusCode int, body []byte, logger logger.Logger) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if len(body) != 0 {
		if _, err := w.Write(body); err != nil {
			logger.Errorw("failed to write response body", "error", err)
		}
	}
}

func ObjectResponse(obj any, message string, logger logger.Logger) *bytes.Buffer {
	response := struct {
		Date    time.Time `json:"date,omitempty"`
		Result  any       `json:"result,omitempty"`
		Message string    `json:"message,omitempty"`
	}{
		Message: message,
		Result:  obj,
		Date:    time.Now().UTC(),
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(response); err != nil {
		logger.Errorw("failed to encode response object to JSON", "error", err)
	}

	return body
}

func HandleError(w http.ResponseWriter, logger logger.Logger, status int, msg string, err error) {
	if err != nil {
		logger.Errorw("operation failed",
			"message", msg,
			"error", err.Error(),
			"status", status,
		)
		response := ObjectResponse(nil, msg+": "+err.Error(), logger)
		ResponseReturn(w, status, response.Bytes(), logger)
	} else {
		logger.Warnw("operation returned warning",
			"message", msg,
			"status", status,
		)
		response := ObjectResponse(nil, msg, logger)
		ResponseReturn(w, status, response.Bytes(), logger)
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
