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

	cmd := buildUpdateCommand(in)

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

// buildUpdateCommand constructs an UpdateRecordCommand from GraphQL input.
func buildUpdateCommand(in gmodel.UpdateRecordInput) input.UpdateRecordCommand {
	cmd := input.UpdateRecordCommand{
		Title:       in.Title,
		Description: in.Description,
		Value:       in.Value,
		Source:      in.Source,
		Timezone:    in.Timezone,
		Status:      in.Status,
	}

	parseUint64Field(in.CategoryID, &cmd.CategoryID)
	parseUint64Field(in.TagID, &cmd.TagID)
	parseTimeField(in.EventTime, &cmd.EventTime)
	parseTimeField(in.RecordedAt, &cmd.RecordedAt)

	if in.DurationSeconds != nil {
		d := int(*in.DurationSeconds)
		cmd.DurationSecs = &d
	}

	return cmd
}

// parseUint64Field parses an optional string field to uint64.
func parseUint64Field(src *string, dst **uint64) {
	if src != nil && *src != "" {
		if v, err := strconv.ParseUint(*src, 10, 64); err == nil {
			*dst = &v
		}
	}
}

// parseTimeField parses an optional RFC3339 string field to time.Time.
func parseTimeField(src *string, dst **time.Time) {
	if src != nil && *src != "" {
		if d, err := time.Parse(time.RFC3339, *src); err == nil {
			*dst = &d
		}
	}
}
