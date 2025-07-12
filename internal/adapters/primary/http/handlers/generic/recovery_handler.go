package generic

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/generic/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
)

// RecoveryHandler handles panics with a standardized error body.
func (h *Handler) RecoveryHandler(w http.ResponseWriter, r *http.Request, recovered interface{}, errorID string) {
	ctx := r.Context()

	reqID := r.Header.Get(commonkeys.XRequestID)
	stack := string(debug.Stack())
	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	_, span := otel.Tracer(constants.TracerGenericHandler).
		Start(ctx, constants.TracerRecoveryHandler)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.URLPath, r.URL.Path),
		attribute.String(commonkeys.RequestID, reqID),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
		attribute.String(tracingkeys.ErrorID, errorID),
	)
	span.RecordError(fmt.Errorf(constants.StacktraceFormat, recovered, stack))
	span.SetStatus(codes.Error, constants.MsgRecoveredFromPanic)

	h.Logger.ErrorwCtx(ctx, constants.MsgRecoveryHandlerFired,
		tracingkeys.RecoveredKey, recovered,
		tracingkeys.StacktraceKey, stack,
		tracingkeys.ErrorID, errorID,
		commonkeys.URLPath, r.URL.Path,
		commonkeys.RequestID, reqID,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	panicErr := fmt.Errorf(constants.RecoveredFormat, constants.MsgRecoveredFromPanic, errorID)
	httpresponse.WriteError(
		w,
		panicErr,
		constants.MsgRecoveredFromPanic,
		h.Logger,
	)
}
