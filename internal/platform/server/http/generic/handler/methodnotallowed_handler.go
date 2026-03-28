// Package handler generic provides common HTTP controllers for the application.
package handler

import (
	"net/http"
	"time"

	"github.com/lechitz/aion-api/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
)

// MethodNotAllowedHandler handles 405 Method Not Allowed responses with a standardized error body.
func (h *Handler) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	reqID, _ := r.Context().Value(ctxkeys.RequestID).(string)
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
