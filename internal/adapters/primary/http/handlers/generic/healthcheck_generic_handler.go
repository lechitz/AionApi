package generic

import (
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/generic/dto"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/generic/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/httpresponse"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// HealthCheck responds to health check requests with service metadata and status.
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	_, span := otel.Tracer(constants.TracerGenericHandler).
		Start(r.Context(), constants.TracerHealthCheckHandler)
	defer span.End()

	reqID := r.Header.Get(commonkeys.XRequestID)

	span.AddEvent(constants.EventHealthCheck,
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
	span.SetStatus(codes.Ok, constants.StatusHealthCheckOK)

	payload := dto.HealthCheckResponse{
		Timestamp: time.Now().UTC(),
		Status:    constants.HealthStatusHealthy,
		Name:      h.GeneralConfig.Name,
		Env:       h.GeneralConfig.Env,
		Version:   h.GeneralConfig.Version,
	}

	h.Logger.Infow(constants.MsgServiceIsHealthy,
		commonkeys.APIName, h.GeneralConfig.Name,
		commonkeys.AppEnv, h.GeneralConfig.Env,
		commonkeys.AppVersion, h.GeneralConfig.Version,
		commonkeys.RequestID, reqID,
	)

	httpresponse.WriteSuccess(w, http.StatusOK, payload, constants.MsgServiceIsHealthy)
}
