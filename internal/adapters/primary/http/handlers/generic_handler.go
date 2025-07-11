// Package handlers provide common HTTP handlers for the application.
package handlers

import (
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/dto"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/response"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/shared/commonkeys"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Generic provides common HTTP handlers for the application.
type Generic struct {
	Logger        output.ContextLogger
	GeneralConfig config.GeneralConfig
}

// NewGeneric initializes and returns a new Generic instance with a Logger dependency.
func NewGeneric(logger output.ContextLogger, general config.GeneralConfig) *Generic {
	return &Generic{
		Logger:        logger,
		GeneralConfig: general,
	}
}

// HealthCheckHandler handles health check requests.
func (h *Generic) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer(constants.TracerGenericHandler).
		Start(r.Context(), constants.TracerHealthCheckHandler)
	defer span.End()

	reqID := r.Header.Get(commonkeys.XRequestID)

	span.SetAttributes(
		attribute.String(commonkeys.APIName, h.GeneralConfig.Name),
		attribute.String(commonkeys.AppEnv, h.GeneralConfig.Env),
		attribute.String(commonkeys.AppVersion, h.GeneralConfig.Version),
		attribute.String(commonkeys.RequestID, reqID),
	)

	span.SetStatus(codes.Ok, constants.MsgSpanHealthCheckOK)

	payload := dto.HealthCheckResponse{
		Status:    constants.HealthStatusHealthy,
		Name:      h.GeneralConfig.Name,
		Env:       h.GeneralConfig.Env,
		Version:   h.GeneralConfig.Version,
		Timestamp: time.Now().UTC(),
	}

	h.Logger.Infow(constants.MsgServiceIsHealthy,
		commonkeys.APIName, h.GeneralConfig.Name,
		commonkeys.AppEnv, h.GeneralConfig.Env,
		commonkeys.AppVersion, h.GeneralConfig.Version,
		commonkeys.RequestID, reqID,
	)

	body := response.ObjectResponse(payload, constants.MsgServiceIsHealthy, h.Logger)

	response.Return(w, http.StatusOK, body.Bytes(), h.Logger)
}

// NotFoundHandler handles 404 errors.
func (h *Generic) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get(commonkeys.XRequestID)

	h.Logger.Infow(constants.MsgResourceNotFound,
		commonkeys.URLPath, r.URL.Path,
		commonkeys.RequestID, reqID,
	)

	payload := dto.ErrorResponse{
		Status:    constants.HealthStatusUnhealthy,
		Message:   constants.MsgResourceNotFound,
		Timestamp: time.Now().UTC(),
	}

	body := response.ObjectResponse(payload, constants.MsgResourceNotFound, h.Logger)

	response.Return(w, http.StatusNotFound, body.Bytes(), h.Logger)
}
