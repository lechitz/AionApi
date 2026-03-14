package controller

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
)

// ListMetricDefinitions returns active metric definitions with multi-tag support.
func (c *controller) ListMetricDefinitions(ctx context.Context, userID uint64) ([]*model.MetricDefinition, error) {
	defs, err := c.RecordService.ListMetricDefinitions(ctx, userID)
	if err != nil {
		return nil, err
	}

	out := make([]*model.MetricDefinition, 0, len(defs))
	for _, def := range defs {
		item := &model.MetricDefinition{
			ID:          strconv.FormatUint(def.ID, 10),
			MetricKey:   def.MetricKey,
			DisplayName: def.DisplayName,
			TagID:       strconv.FormatUint(def.TagID, 10),
			TagIds:      formatIDs(def.TagIDs),
			ValueSource: def.ValueSource,
			Aggregation: def.Aggregation,
			Unit:        def.Unit,
			GoalDefault: def.GoalDefault,
			IsActive:    def.IsActive,
		}
		if def.CategoryID != nil {
			catID := strconv.FormatUint(*def.CategoryID, 10)
			item.CategoryID = &catID
		}
		out = append(out, item)
	}
	return out, nil
}
