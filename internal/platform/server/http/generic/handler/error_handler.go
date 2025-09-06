// Package generic provides common HTTP controllers for the application.
package handler

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
)

// ErrorHandler handles errors with a standardized error body.
func (h *Handler) ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	reqID := ctx.Value(ctxkeys.RequestID)
	reqIDStr, _ := reqID.(string)
	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	_, span := otel.Tracer(TracerGenericHandler).
		Start(ctx, TracerErrorHandler)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.URLPath, r.URL.Path),
		attribute.String(commonkeys.RequestID, reqIDStr),
		attribute.String(commonkeys.Method, r.Method),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)
	span.RecordError(err)
	span.SetStatus(codes.Error, MsgInternalServerError)

	h.Logger.Errorw(MsgInternalServerError,
		commonkeys.URLPath, r.URL.Path,
		commonkeys.RequestID, reqIDStr,
		commonkeys.Error, err,
	)

	httpresponse.WriteError(
		w,
		ErrInternalServer,
		MsgInternalServerError,
		h.Logger,
	)
}
