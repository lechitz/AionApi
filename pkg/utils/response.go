package utils

import (
	"bytes"
	"encoding/json"
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
		Date:    time.Now(),
	}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(response)
	return body
}
