// Package generic provides common HTTP controllers for the application.
package handler

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/shared/httpresponse"
)

// MethodNotAllowedHandler handles 405 Method Not Allowed responses with a standardized error body.
func (h *Handler) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get(commonkeys.XRequestID)
	h.Logger.Warnw(MsgMethodNotAllowed,
		commonkeys.Method, r.Method,
		commonkeys.URLPath, r.URL.Path,
		commonkeys.RequestID, reqID,
	)
	httpresponse.WriteError(
		w,
		ErrMethodNotAllowed,
		MsgMethodNotAllowed,
		h.Logger,
	)
}
