package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lechitz/aion-api/internal/record/core/domain"
)

type recordProjectionRow struct {
	RecordID           uint64          `gorm:"column:record_id"`
	UserID             uint64          `gorm:"column:user_id"`
	TagID              uint64          `gorm:"column:tag_id"`
	Description        *string         `gorm:"column:description"`
	EventTimeUTC       time.Time       `gorm:"column:event_time_utc"`
	RecordedAtUTC      *time.Time      `gorm:"column:recorded_at_utc"`
	Status             *string         `gorm:"column:status"`
	Timezone           *string         `gorm:"column:timezone"`
	DurationSeconds    *int            `gorm:"column:duration_seconds"`
	Value              *float64        `gorm:"column:value"`
	Source             *string         `gorm:"column:source"`
	LastEventID        string          `gorm:"column:last_event_id"`
	LastEventType      string          `gorm:"column:last_event_type"`
	LastEventVersion   string          `gorm:"column:last_event_version"`
	LastTraceID        *string         `gorm:"column:last_trace_id"`
	LastRequestID      *string         `gorm:"column:last_request_id"`
	LastKafkaTopic     string          `gorm:"column:last_kafka_topic"`
	LastKafkaPartition int             `gorm:"column:last_kafka_partition"`
	LastKafkaOffset    int64           `gorm:"column:last_kafka_offset"`
	LastConsumedAtUTC  time.Time       `gorm:"column:last_consumed_at_utc"`
	PayloadJSON        json.RawMessage `gorm:"column:payload_json"`
	CreatedAtUTC       time.Time       `gorm:"column:created_at_utc"`
	UpdatedAtUTC       time.Time       `gorm:"column:updated_at_utc"`
}

func toRecordProjection(row recordProjectionRow) domain.RecordProjection {
	return domain.RecordProjection{
		RecordID:           row.RecordID,
		UserID:             row.UserID,
		TagID:              row.TagID,
		Description:        row.Description,
		EventTimeUTC:       row.EventTimeUTC,
		RecordedAtUTC:      row.RecordedAtUTC,
		Status:             row.Status,
		Timezone:           row.Timezone,
		DurationSeconds:    row.DurationSeconds,
		Value:              row.Value,
		Source:             row.Source,
		LastEventID:        row.LastEventID,
		LastEventType:      row.LastEventType,
		LastEventVersion:   row.LastEventVersion,
		LastTraceID:        row.LastTraceID,
		LastRequestID:      row.LastRequestID,
		LastKafkaTopic:     row.LastKafkaTopic,
		LastKafkaPartition: row.LastKafkaPartition,
		LastKafkaOffset:    row.LastKafkaOffset,
		LastConsumedAtUTC:  row.LastConsumedAtUTC,
		PayloadJSON:        []byte(row.PayloadJSON),
		CreatedAtUTC:       row.CreatedAtUTC,
		UpdatedAtUTC:       row.UpdatedAtUTC,
	}
}

// GetProjectedByID returns one derived record projection owned by aion_derived.
func (r *RecordRepository) GetProjectedByID(ctx context.Context, userID uint64, recordID uint64) (domain.RecordProjection, error) {
	var row recordProjectionRow
	if err := r.db.WithContext(ctx).
		Raw(`
			SELECT
				record_id,
				user_id,
				tag_id,
				description,
				event_time_utc,
				recorded_at_utc,
				status,
				timezone,
				duration_seconds,
				value,
				source,
				last_event_id,
				last_event_type,
				last_event_version,
				last_trace_id,
				last_request_id,
				last_kafka_topic,
				last_kafka_partition,
				last_kafka_offset,
				last_consumed_at_utc,
				payload_json,
				created_at_utc,
				updated_at_utc
			FROM aion_derived.record_projection_v1
			WHERE user_id = ? AND record_id = ?
			LIMIT 1
		`, userID, recordID).
		Scan(&row).Error(); err != nil {
		return domain.RecordProjection{}, fmt.Errorf("get projected record: %w", err)
	}

	if row.RecordID == 0 {
		return domain.RecordProjection{}, errors.New("get projected record: record not found")
	}

	return toRecordProjection(row), nil
}

// ListProjectedLatest returns the latest derived projections ordered by last consumed offset.
func (r *RecordRepository) ListProjectedLatest(ctx context.Context, userID uint64, limit int) ([]domain.RecordProjection, error) {
	var rows []recordProjectionRow
	if err := r.db.WithContext(ctx).
		Raw(`
			SELECT
				record_id,
				user_id,
				tag_id,
				description,
				event_time_utc,
				recorded_at_utc,
				status,
				timezone,
				duration_seconds,
				value,
				source,
				last_event_id,
				last_event_type,
				last_event_version,
				last_trace_id,
				last_request_id,
				last_kafka_topic,
				last_kafka_partition,
				last_kafka_offset,
				last_consumed_at_utc,
				payload_json,
				created_at_utc,
				updated_at_utc
			FROM aion_derived.record_projection_v1
			WHERE user_id = ?
			ORDER BY last_consumed_at_utc DESC, record_id DESC
			LIMIT ?
		`, userID, limit).
		Scan(&rows).Error(); err != nil {
		return nil, fmt.Errorf("list projected records: %w", err)
	}

	out := make([]domain.RecordProjection, len(rows))
	for i := range rows {
		out[i] = toRecordProjection(rows[i])
	}
	return out, nil
}

// ListProjectedPage returns derived projections ordered by event time desc with the same cursor contract as records().
func (r *RecordRepository) ListProjectedPage(ctx context.Context, userID uint64, limit int, afterEventTime *string, afterID *int64) ([]domain.RecordProjection, error) {
	var rows []recordProjectionRow
	args := []any{userID}
	cursorClause := ""
	if afterEventTime != nil && afterID != nil {
		cursorClause = `
			AND (event_time_utc < ? OR (event_time_utc = ? AND record_id < ?))
		`
		args = append(args, *afterEventTime, *afterEventTime, *afterID)
	}
	args = append(args, limit)

	if err := r.db.WithContext(ctx).
		Raw(fmt.Sprintf(`
			SELECT
				record_id,
				user_id,
				tag_id,
				description,
				event_time_utc,
				recorded_at_utc,
				status,
				timezone,
				duration_seconds,
				value,
				source,
				last_event_id,
				last_event_type,
				last_event_version,
				last_trace_id,
				last_request_id,
				last_kafka_topic,
				last_kafka_partition,
				last_kafka_offset,
				last_consumed_at_utc,
				payload_json,
				created_at_utc,
				updated_at_utc
			FROM aion_derived.record_projection_v1
			WHERE user_id = ?
			%s
			ORDER BY event_time_utc DESC, record_id DESC
			LIMIT ?
		`, cursorClause), args...).
		Scan(&rows).Error(); err != nil {
		return nil, fmt.Errorf("list projected page: %w", err)
	}

	out := make([]domain.RecordProjection, len(rows))
	for i := range rows {
		out[i] = toRecordProjection(rows[i])
	}
	return out, nil
}
