package handler

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// NotFoundHandler handles 404 resource not found responses with a standardized error body.
func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get(commonkeys.XRequestID)

	h.Logger.Warnw(MsgResourceNotFound,
		commonkeys.URLPath, r.URL.Path,
		commonkeys.RequestID, reqID,
	)

	httpresponse.WriteError(
		w,
		ErrResourceNotFound,
		MsgResourceNotFound,
		h.Logger,
	)
}
