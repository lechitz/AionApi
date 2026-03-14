package usecase

import (
	"context"
	"encoding/json"
	"strconv"

	eventoutboxdomain "github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
	eventoutboxinput "github.com/lechitz/AionApi/internal/eventoutbox/core/ports/input"
	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

func (s *Service) enqueueRecordOutboxEventWithService(ctx context.Context, outboxService eventoutboxinput.Service, eventType string, record domain.Record) {
	if outboxService == nil {
		return
	}

	traceID, _ := ctx.Value(ctxkeys.TraceID).(string)
	requestID, _ := ctx.Value(ctxkeys.RequestID).(string)
	payloadJSON, err := json.Marshal(map[string]any{
		"record_id":        record.ID,
		"user_id":          record.UserID,
		"tag_id":           record.TagID,
		"event_time_utc":   record.EventTime.UTC().Format("2006-01-02T15:04:05.000000Z07:00"),
		"recorded_at_utc":  record.RecordedAt,
		"status":           record.Status,
		"timezone":         record.Timezone,
		"duration_seconds": record.DurationSecs,
		"value":            record.Value,
		"source":           record.Source,
		"description":      record.Description,
	})
	if err != nil {
		s.Logger.WarnwCtx(ctx, LogFailedToMarshalRecordEventPayload,
			commonkeys.Error, err,
			commonkeys.RecordID, record.ID,
			commonkeys.UserID, record.UserID,
			"event_type", eventType,
		)
		return
	}

	event := eventoutboxdomain.Event{
		AggregateType: RecordAggregateType,
		AggregateID:   strconv.FormatUint(record.ID, 10),
		EventType:     eventType,
		EventVersion:  RecordEventVersionV1,
		TraceID:       traceID,
		RequestID:     requestID,
		PayloadJSON:   payloadJSON,
	}

	if err := outboxService.Enqueue(ctx, event); err != nil {
		s.Logger.WarnwCtx(ctx, LogFailedToEnqueueRecordEvent,
			commonkeys.Error, err,
			commonkeys.RecordID, record.ID,
			commonkeys.UserID, record.UserID,
			"event_type", eventType,
		)
	}
}
