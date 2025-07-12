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
func (h *Handler) RecoveryHandler(w http.ResponseWriter, r *http.Request, recovered interface{}) {
	reqID := r.Header.Get(commonkeys.XRequestID)
	stack := string(debug.Stack())
	ip := r.RemoteAddr
	userAgent := r.UserAgent()

	_, span := otel.Tracer(constants.TracerGenericHandler).
		Start(r.Context(), constants.TracerRecoveryHandler)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.URLPath, r.URL.Path),
		attribute.String(commonkeys.RequestID, reqID),
		attribute.String(tracingkeys.RequestIPKey, ip),
		attribute.String(tracingkeys.RequestUserAgentKey, userAgent),
	)
	span.RecordError(fmt.Errorf(constants.StacktraceFormat, recovered, stack))
	span.SetStatus(codes.Error, constants.MsgRecoveredFromPanic)

	h.Logger.Errorw(constants.MsgRecoveredFromPanic,
		tracingkeys.RecoveredKey, recovered,
		tracingkeys.StacktraceKey, stack,
		commonkeys.URLPath, r.URL.Path,
		commonkeys.RequestID, reqID,
		tracingkeys.RequestIPKey, ip,
		tracingkeys.RequestUserAgentKey, userAgent,
	)

	httpresponse.WriteError(
		w,
		constants.ErrRecoveredFromPanic,
		constants.MsgRecoveredFromPanic,
		h.Logger,
	)
}
