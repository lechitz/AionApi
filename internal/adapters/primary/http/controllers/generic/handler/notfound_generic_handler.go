package handler

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/generic/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
)

// NotFoundHandler handles 404 resource not found responses with a standardized error body.
func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get(commonkeys.XRequestID)

	h.Logger.Warnw(constants.MsgResourceNotFound,
		commonkeys.URLPath, r.URL.Path,
		commonkeys.RequestID, reqID,
	)

	httpresponse.WriteError(
		w,
		constants.ErrResourceNotFound,
		constants.MsgResourceNotFound,
		h.Logger,
	)
}
