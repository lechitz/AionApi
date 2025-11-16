// Package controller provides mapping helpers between GraphQL models and core commands/domain for the Record context.
package controller

import (
	"strconv"
	"time"

	gmodel "github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
)

// toModelOut converts a domain.Record to a GraphQL model.Record.
func toModelOut(t domain.Record) *gmodel.Record {
	out := &gmodel.Record{
		ID:         strconv.FormatUint(t.ID, 10),
		UserID:     strconv.FormatUint(t.UserID, 10),
		CategoryID: strconv.FormatUint(t.CategoryID, 10),
		Title:      t.Title,
		EventTime:  t.EventTime.UTC().Format(time.RFC3339),
		CreatedAt:  t.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:  t.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if t.Description != nil {
		out.Description = t.Description
	}
	if t.RecordedAt != nil {
		r := t.RecordedAt.UTC().Format(time.RFC3339)
		out.RecordedAt = &r
	}
	if t.DurationSecs != nil {
		d := int32(*t.DurationSecs)
		out.DurationSeconds = &d
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
	// TagID is required
	outTag := strconv.FormatUint(t.TagID, 10)
	out.TagID = outTag
	return out
}

// toCreateCommand converts a GraphQL CreateRecordInput into an input.CreateRecordCommand.
func toCreateCommand(in gmodel.CreateRecordInput, userID uint64) input.CreateRecordCommand {
	uid := userID

	var categoryID uint64
	if v, err := strconv.ParseUint(in.CategoryID, 10, 64); err == nil {
		categoryID = v
	}

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
		CategoryID:   categoryID,
		Title:        in.Title,
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
