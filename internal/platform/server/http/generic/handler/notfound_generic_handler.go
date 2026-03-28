package handler

import (
	"net/http"
	"time"

	"github.com/lechitz/aion-api/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
)

// NotFoundHandler handles 404 resource not found responses with a standardized error body.
func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	reqID, _ := r.Context().Value(ctxkeys.RequestID).(string)

	h.Logger.Warnw(MsgResourceNotFound,
		commonkeys.URLPath, r.URL.Path,
		commonkeys.RequestID, reqID,
	)

	body := httpresponse.ResponseBody{
		Date:    time.Now().UTC(),
		Error:   MsgResourceNotFound,
		Details: ErrResourceNotFound.Error(),
		Code:    http.StatusNotFound,
	}
	httpresponse.WriteJSON(w, http.StatusNotFound, body)
}
