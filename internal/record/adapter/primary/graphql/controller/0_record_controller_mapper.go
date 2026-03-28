// Package controller provides mapping helpers between GraphQL models and core commands/domain for the Record context.
package controller

import (
	"math"
	"strconv"
	"time"

	gmodel "github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/ports/input"
)

// toModelOut converts a domain.Record to a GraphQL model.Record.
func toModelOut(t domain.Record) *gmodel.Record {
	out := &gmodel.Record{
		ID:        strconv.FormatUint(t.ID, 10),
		UserID:    strconv.FormatUint(t.UserID, 10),
		TagID:     strconv.FormatUint(t.TagID, 10),
		EventTime: t.EventTime.UTC().Format(time.RFC3339),
		CreatedAt: t.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if t.Description != nil {
		out.Description = t.Description
	}
	if t.RecordedAt != nil {
		r := t.RecordedAt.UTC().Format(time.RFC3339)
		out.RecordedAt = &r
	}
	if t.DurationSecs != nil {
		dv := safeRecordIntToInt32(*t.DurationSecs)
		out.DurationSeconds = &dv
	}
	if t.Value != nil {
		out.Value = t.Value
	}
	if t.Source != nil {
		out.Source = t.Source
	}
	if t.Timezone != nil {
		out.Timezone = t.Timezone
	}
	if t.Status != nil {
		out.Status = t.Status
	}
	return out
}

// toModelOutSlice converts a slice of domain.Record to a slice of GraphQL model.Record pointers.
func toModelOutSlice(records []domain.Record) []*gmodel.Record {
	result := make([]*gmodel.Record, len(records))
	for i, rec := range records {
		result[i] = toModelOut(rec)
	}
	return result
}

func toProjectedModelOut(t domain.RecordProjection) *gmodel.RecordProjection {
	out := &gmodel.RecordProjection{
		RecordID:           strconv.FormatUint(t.RecordID, 10),
		UserID:             strconv.FormatUint(t.UserID, 10),
		TagID:              strconv.FormatUint(t.TagID, 10),
		EventTimeUtc:       t.EventTimeUTC.UTC().Format(time.RFC3339),
		LastEventID:        t.LastEventID,
		LastEventType:      t.LastEventType,
		LastEventVersion:   t.LastEventVersion,
		LastKafkaTopic:     t.LastKafkaTopic,
		LastKafkaPartition: safeRecordIntToInt32(t.LastKafkaPartition),
		LastKafkaOffset:    strconv.FormatInt(t.LastKafkaOffset, 10),
		LastConsumedAtUtc:  t.LastConsumedAtUTC.UTC().Format(time.RFC3339),
		PayloadJSON:        string(t.PayloadJSON),
		CreatedAtUtc:       t.CreatedAtUTC.UTC().Format(time.RFC3339),
		UpdatedAtUtc:       t.UpdatedAtUTC.UTC().Format(time.RFC3339),
	}
	if t.Description != nil {
		out.Description = t.Description
	}
	if t.RecordedAtUTC != nil {
		v := t.RecordedAtUTC.UTC().Format(time.RFC3339)
		out.RecordedAtUtc = &v
	}
	if t.Status != nil {
		out.Status = t.Status
	}
	if t.Timezone != nil {
		out.Timezone = t.Timezone
	}
	if t.DurationSeconds != nil {
		v := safeRecordIntToInt32(*t.DurationSeconds)
		out.DurationSeconds = &v
	}
	if t.Value != nil {
		out.Value = t.Value
	}
	if t.Source != nil {
		out.Source = t.Source
	}
	if t.LastTraceID != nil {
		out.LastTraceID = t.LastTraceID
	}
	if t.LastRequestID != nil {
		out.LastRequestID = t.LastRequestID
	}
	return out
}

func toProjectedModelOutSlice(items []domain.RecordProjection) []*gmodel.RecordProjection {
	result := make([]*gmodel.RecordProjection, len(items))
	for i, item := range items {
		result[i] = toProjectedModelOut(item)
	}
	return result
}

// toCreateCommand converts a GraphQL CreateRecordInput into an input.CreateRecordCommand.
func toCreateCommand(in gmodel.CreateRecordInput, userID uint64) input.CreateRecordCommand {
	uid := userID

	// parse eventTime
	var eventTime time.Time
	if in.EventTime != nil && *in.EventTime != "" {
		if d, err := time.Parse(time.RFC3339, *in.EventTime); err == nil {
			eventTime = d
		}
	}

	// parse recordedAt
	var recordedAt *time.Time
	if in.RecordedAt != nil && *in.RecordedAt != "" {
		if d, err := time.Parse(time.RFC3339, *in.RecordedAt); err == nil {
			recordedAt = &d
		}
	}

	// parse single tag id (required) - in.TagID is string
	tagID := uint64(0)
	if in.TagID != "" {
		if v, err := strconv.ParseUint(in.TagID, 10, 64); err == nil {
			tagID = v
		}
	}

	var duration *int
	if in.DurationSeconds != nil {
		d := int(*in.DurationSeconds)
		duration = &d
	}

	return input.CreateRecordCommand{
		UserID:       uid,
		Description:  in.Description,
		TagID:        tagID,
		EventTime:    eventTime,
		RecordedAt:   recordedAt,
		DurationSecs: duration,
		Value:        in.Value,
		Source:       in.Source,
		Timezone:     in.Timezone,
		Status:       in.Status,
	}
}

func safeRecordIntToInt32(value int) int32 {
	if value > math.MaxInt32 {
		return math.MaxInt32
	}
	if value < math.MinInt32 {
		return math.MinInt32
	}
	return int32(value)
}
