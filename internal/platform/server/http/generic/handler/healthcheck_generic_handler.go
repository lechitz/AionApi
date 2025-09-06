package handler

import (
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/adapter/server/http/generic/dto"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// HealthCheck Health Check responds to health check requests with service metadata and status.
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer(TracerGenericHandler).
		Start(r.Context(), TracerHealthCheckHandler)
	defer span.End()

	reqID := r.Header.Get(commonkeys.XRequestID)

	span.AddEvent(EventHealthCheck,
		trace.WithAttributes(
			attribute.String(commonkeys.APIName, h.GeneralConfig.Name),
			attribute.String(commonkeys.AppEnv, h.GeneralConfig.Env),
			attribute.String(commonkeys.AppVersion, h.GeneralConfig.Version),
			attribute.String(commonkeys.RequestID, reqID),
		),
	)

	span.SetAttributes(
		attribute.String(commonkeys.APIName, h.GeneralConfig.Name),
		attribute.String(commonkeys.AppEnv, h.GeneralConfig.Env),
		attribute.String(commonkeys.AppVersion, h.GeneralConfig.Version),
		attribute.String(commonkeys.RequestID, reqID),
	)
	span.SetStatus(codes.Ok, StatusHealthCheckOK)

	payload := dto.HealthCheckResponse{
		Timestamp: time.Now().UTC(),
		Status:    HealthStatusHealthy,
		Name:      h.GeneralConfig.Name,
		Env:       h.GeneralConfig.Env,
		Version:   h.GeneralConfig.Version,
	}

	h.Logger.Infow(MsgServiceIsHealthy,
		commonkeys.APIName, h.GeneralConfig.Name,
		commonkeys.AppEnv, h.GeneralConfig.Env,
		commonkeys.AppVersion, h.GeneralConfig.Version,
		commonkeys.RequestID, reqID,
	)

	httpresponse.WriteSuccess(w, http.StatusOK, payload, MsgServiceIsHealthy)
}
