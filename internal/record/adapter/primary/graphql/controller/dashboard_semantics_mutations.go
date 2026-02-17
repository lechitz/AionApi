package controller

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
)

func (c *controller) UpsertMetricDefinition(ctx context.Context, userID uint64, in model.UpsertMetricDefinitionInput) (*model.MetricDefinition, error) {
	primaryTagID := mustParseID(in.TagID)
	tagIDs := make([]uint64, 0, len(in.TagIds)+1)
	tagIDs = append(tagIDs, primaryTagID)
	for _, rawID := range in.TagIds {
		tagIDs = append(tagIDs, mustParseID(rawID))
	}

	cmd := input.UpsertMetricDefinitionCommand{
		MetricKey:   in.MetricKey,
		DisplayName: in.DisplayName,
		TagID:       primaryTagID,
		TagIDs:      tagIDs,
		ValueSource: toPtrValue(in.ValueSource),
		Aggregation: toPtrValue(in.Aggregation),
		Unit:        toPtrValue(in.Unit),
		GoalDefault: in.GoalDefault,
		IsActive:    in.IsActive,
	}
	if in.ID != nil {
		id := mustParseID(*in.ID)
		cmd.ID = &id
	}
	if in.CategoryID != nil {
		cid := mustParseID(*in.CategoryID)
		cmd.CategoryID = &cid
	}

	out, err := c.RecordService.UpsertMetricDefinition(ctx, userID, cmd)
	if err != nil {
		return nil, err
	}

	result := &model.MetricDefinition{
		ID:          strconv.FormatUint(out.ID, 10),
		MetricKey:   out.MetricKey,
		DisplayName: out.DisplayName,
		TagID:       strconv.FormatUint(out.TagID, 10),
		TagIds:      formatIDs(out.TagIDs),
		ValueSource: out.ValueSource,
		Aggregation: out.Aggregation,
		Unit:        out.Unit,
		GoalDefault: out.GoalDefault,
		IsActive:    out.IsActive,
	}
	if out.CategoryID != nil {
		value := strconv.FormatUint(*out.CategoryID, 10)
		result.CategoryID = &value
	}
	return result, nil
}

func (c *controller) UpsertGoalTemplate(ctx context.Context, userID uint64, in model.UpsertGoalTemplateInput) (*model.GoalTemplate, error) {
	cmd := input.UpsertGoalTemplateCommand{
		MetricKey:   in.MetricKey,
		Title:       in.Title,
		TargetValue: in.TargetValue,
		Comparison:  toPtrValue(in.Comparison),
		Period:      toPtrValue(in.Period),
		IsActive:    in.IsActive,
	}
	if in.ID != nil {
		id := mustParseID(*in.ID)
		cmd.ID = &id
	}

	out, err := c.RecordService.UpsertGoalTemplate(ctx, userID, cmd)
	if err != nil {
		return nil, err
	}

	return &model.GoalTemplate{
		ID:          strconv.FormatUint(out.ID, 10),
		MetricKey:   out.MetricKey,
		Title:       out.Title,
		TargetValue: out.TargetValue,
		Comparison:  out.Comparison,
		Period:      out.Period,
		IsActive:    out.IsActive,
	}, nil
}

func (c *controller) DeleteGoalTemplate(ctx context.Context, userID uint64, id uint64) error {
	return c.RecordService.DeleteGoalTemplate(ctx, userID, id)
}

func mustParseID(id string) uint64 {
	value, _ := strconv.ParseUint(id, 10, 64)
	return value
}

func toPtrValue(in *string) string {
	if in == nil {
		return ""
	}
	return *in
}

func formatIDs(values []uint64) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == 0 {
			continue
		}
		out = append(out, strconv.FormatUint(value, 10))
	}
	return out
}
