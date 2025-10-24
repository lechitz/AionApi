// Package handler generic provides common HTTP controllers for the application.
package handler

import (
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// MethodNotAllowedHandler handles 405 Method Not Allowed responses with a standardized error body.
func (h *Handler) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get(commonkeys.XRequestID)
	h.Logger.Warnw(MsgMethodNotAllowed,
		commonkeys.Method, r.Method,
		commonkeys.URLPath, r.URL.Path,
		commonkeys.RequestID, reqID,
	)

	body := httpresponse.ResponseBody{
		Date:    time.Now().UTC(),
		Error:   MsgMethodNotAllowed,
		Details: ErrMethodNotAllowed.Error(),
		Code:    http.StatusMethodNotAllowed,
	}
	httpresponse.WriteJSON(w, http.StatusMethodNotAllowed, body)
}
