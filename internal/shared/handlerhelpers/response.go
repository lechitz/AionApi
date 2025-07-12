// Package handlerhelpers provides common response handling functions.
package handlerhelpers

import (
	"context"
	"net/http"

	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/tracingkeys"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// WriteAuthError handles authentication error response.
func WriteAuthError(ctx context.Context, w http.ResponseWriter, span trace.Span, logger output.ContextLogger) {
	err := sharederrors.ErrMissingUserID()
	span.RecordError(err)
	span.SetStatus(codes.Error, "missing user id")
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusUnauthorized))
	logger.ErrorwCtx(ctx, "missing user id", commonkeys.Error, err.Error())
	httpresponse.WriteAuthError(w, err, logger)
}

// WriteDecodeError handles request decode error response.
func WriteDecodeError(ctx context.Context, w http.ResponseWriter, span trace.Span, err error, logger output.ContextLogger) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, http.StatusBadRequest))
	logger.ErrorwCtx(ctx, "decode error", commonkeys.Error, err.Error())
	httpresponse.WriteDecodeError(w, err, logger)
}

// WriteDomainError handles domain error response.
func WriteDomainError(ctx context.Context, w http.ResponseWriter, span trace.Span, err error, msg string, logger output.ContextLogger) {
	statusCode := sharederrors.MapErrorToHTTPStatus(err)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	span.SetAttributes(attribute.Int(tracingkeys.HTTPStatusCodeKey, statusCode))
	logger.ErrorwCtx(ctx, "domain error", commonkeys.Error, err.Error())
	httpresponse.WriteError(w, err, msg, logger)
}
