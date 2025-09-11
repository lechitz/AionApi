package handler

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/httpresponse"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// RecoveryHandler handles panics with a standardized error body.
func (h *Handler) RecoveryHandler(w http.ResponseWriter, r *http.Request, recovered interface{}, errorID string) {
	ctx := r.Context()

	reqID := r.Header.Get(commonkeys.XRequestID)
	stack := string(debug.Stack())
	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	_, span := otel.Tracer(TracerGenericHandler).
		Start(ctx, TracerRecoveryHandler)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.URLPath, r.URL.Path),
		attribute.String(commonkeys.RequestID, reqID),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		attribute.String(tracingkeys.ErrorID, errorID),
	)
	span.RecordError(fmt.Errorf(StacktraceFormat, recovered, stack))
	span.SetStatus(codes.Error, MsgRecoveredFromPanic)

	h.Logger.ErrorwCtx(ctx, MsgRecoveryHandlerFired,
		tracingkeys.RecoveredKey, recovered,
		tracingkeys.StacktraceKey, stack,
		tracingkeys.ErrorID, errorID,
		commonkeys.URLPath, r.URL.Path,
		commonkeys.RequestID, reqID,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	panicErr := fmt.Errorf(RecoveredFormat, MsgRecoveredFromPanic, errorID)
	httpresponse.WriteError(
		w,
		panicErr,
		MsgRecoveredFromPanic,
		h.Logger,
	)
}
