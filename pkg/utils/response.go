package utils

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"time"
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

func HandleError(w http.ResponseWriter, logger *zap.SugaredLogger, status int, msg string, err error) {
	if err != nil {
		logger.Errorw(msg, "error", err.Error())
		response := ObjectResponse(nil, msg+": "+err.Error())
		ResponseReturn(w, status, response.Bytes())
	} else {
		logger.Errorw(msg)
		response := ObjectResponse(nil, msg)
		ResponseReturn(w, status, response.Bytes())
	}
}

func HandleCriticalError(loggerSugar *zap.SugaredLogger, message string, err error) {
	if err != nil {
		loggerSugar.Fatalw(message, "error", err.Error())
	} else {
		loggerSugar.Fatal(message)
	}
}
