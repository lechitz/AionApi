package controller

import (
	"context"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
)

// DashboardSnapshot returns aggregated deterministic metrics and goals for the day.
func (c *controller) DashboardSnapshot(ctx context.Context, userID uint64, date string, timezone *string) (*model.DashboardSnapshot, error) {
	targetDate, err := parseDateOrDefault(date)
	if err != nil {
		return nil, err
	}

	tz := ""
	if timezone != nil {
		tz = *timezone
	}

	out, err := c.RecordService.DashboardSnapshot(ctx, userID, input.DashboardSnapshotQuery{
		Date:     targetDate,
		Timezone: tz,
	})
	if err != nil {
		return nil, err
	}

	metrics := make([]*model.DashboardMetric, 0, len(out.Metrics))
	for _, metric := range out.Metrics {
		item := &model.DashboardMetric{
			MetricKey:   metric.MetricKey,
			Label:       metric.Label,
			Value:       metric.Value,
			Unit:        metric.Unit,
			ProgressPct: metric.ProgressPct,
			Status:      metric.Status,
		}
		if metric.Target != nil {
			item.Target = metric.Target
		}
		metrics = append(metrics, item)
	}

	goals := make([]*model.DashboardGoal, 0, len(out.Goals))
	for _, goal := range out.Goals {
		goals = append(goals, &model.DashboardGoal{
			GoalID:       strconv.FormatUint(goal.GoalID, 10),
			Title:        goal.Title,
			MetricKey:    goal.MetricKey,
			CurrentValue: goal.Current,
			TargetValue:  goal.Target,
			ProgressPct:  goal.ProgressPct,
			Status:       goal.Status,
		})
	}

	return &model.DashboardSnapshot{
		Date:     out.Date.Format("2006-01-02"),
		Timezone: out.Timezone,
		Metrics:  metrics,
		Goals:    goals,
	}, nil
}

func parseDateOrDefault(date string) (time.Time, error) {
	if date == "" {
		return time.Now(), nil
	}
	return time.Parse("2006-01-02", date)
}
