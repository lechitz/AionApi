package controller

import (
	"context"
	"errors"
	"strconv"
	"time"

	gmodel "github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Update updates a record via GraphQL.
func (h *controller) Update(ctx context.Context, in gmodel.UpdateRecordInput, userID uint64) (*gmodel.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "record.update")
	defer span.End()

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound)
		h.Logger.ErrorwCtx(ctx, ErrUserIDNotFound, commonkeys.UserID, userID)
		return nil, errors.New(ErrUserIDNotFound)
	}

	recID, err := strconv.ParseUint(in.ID, 10, 64)
	if err != nil {
		span.SetStatus(codes.Error, ErrInvalidRecordID)
		h.Logger.ErrorwCtx(ctx, ErrInvalidRecordID, "record_id", in.ID, commonkeys.Error, err.Error())
		return nil, errors.New(ErrInvalidRecordID)
	}

	cmd := input.UpdateRecordCommand{}
	if in.Title != nil {
		cmd.Title = in.Title
	}
	if in.Description != nil {
		cmd.Description = in.Description
	}
	if in.CategoryID != nil && *in.CategoryID != "" {
		if v, err := strconv.ParseUint(*in.CategoryID, 10, 64); err == nil {
			cmd.CategoryID = &v
		}
	}
	if in.TagID != nil && *in.TagID != "" {
		if v, err := strconv.ParseUint(*in.TagID, 10, 64); err == nil {
			cmd.TagID = &v
		}
	}
	if in.EventTime != nil && *in.EventTime != "" {
		if d, err := time.Parse(time.RFC3339, *in.EventTime); err == nil {
			cmd.EventTime = &d
		}
	}
	if in.RecordedAt != nil && *in.RecordedAt != "" {
		if d, err := time.Parse(time.RFC3339, *in.RecordedAt); err == nil {
			cmd.RecordedAt = &d
		}
	}
	if in.DurationSeconds != nil {
		d := int(*in.DurationSeconds)
		cmd.DurationSecs = &d
	}
	if in.Value != nil {
		cmd.Value = in.Value
	}
	if in.Source != nil {
		cmd.Source = in.Source
	}
	if in.Timezone != nil {
		cmd.Timezone = in.Timezone
	}
	if in.Status != nil {
		cmd.Status = in.Status
	}

	do, err := h.RecordService.Update(ctx, recID, userID, cmd)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error updating record")
		h.Logger.ErrorwCtx(ctx, "error updating record", commonkeys.Error, err.Error())
		return nil, err
	}

	out := toModelOut(do)
	span.SetAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("record_id", in.ID),
	)
	span.SetStatus(codes.Ok, "updated")
	return out, nil
}
