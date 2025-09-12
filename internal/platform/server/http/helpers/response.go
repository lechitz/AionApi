// Package helpers provides common response handling functions.
package helpers

import (
	"context"
	"net/http"

	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/httpresponse"
	sharederrors2 "github.com/lechitz/AionApi/internal/platform/server/http/helpers/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// WriteAuthError handles authentication error response.
func WriteAuthError(ctx context.Context, w http.ResponseWriter, span trace.Span, logger logger.ContextLogger) {
	err := sharederrors2.ErrMissingUserID()
	span.RecordError(err)
	span.SetStatus(codes.Error, sharederrors2.ErrMsgMissingUserID)
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusUnauthorized))
	logger.ErrorwCtx(ctx, sharederrors2.ErrMsgMissingUserID, commonkeys.Error, err.Error())
	httpresponse.WriteAuthError(w, err, logger)
}

// WriteDecodeError handles request decode error response.
func WriteDecodeError(ctx context.Context, w http.ResponseWriter, span trace.Span, err error, logger logger.ContextLogger) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusBadRequest))
	logger.ErrorwCtx(ctx, "decode error", commonkeys.Error, err.Error())
	httpresponse.WriteDecodeError(w, err, logger)
}

// WriteDomainError handles domain error response.
func WriteDomainError(ctx context.Context, w http.ResponseWriter, span trace.Span, err error, msg string, logger logger.ContextLogger) {
	statusCode := sharederrors2.MapErrorToHTTPStatus(err)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, statusCode))
	logger.ErrorwCtx(ctx, "domain error", commonkeys.Error, err.Error())
	httpresponse.WriteError(w, err, msg, logger)
}
