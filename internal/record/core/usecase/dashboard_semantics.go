package usecase

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
)

// ListMetricDefinitions returns active dashboard metric definitions for the user.
func (s *Service) ListMetricDefinitions(ctx context.Context, userID uint64) ([]domain.MetricDefinition, error) {
	if userID == 0 {
		return nil, ErrUserIDIsRequired
	}
	return s.RecordRepository.ListMetricDefinitions(ctx, userID)
}

// DashboardSnapshot computes deterministic metrics and daily goals for a specific date.
func (s *Service) DashboardSnapshot(ctx context.Context, userID uint64, query input.DashboardSnapshotQuery) (domain.DashboardSnapshot, error) {
	if userID == 0 {
		return domain.DashboardSnapshot{}, ErrUserIDIsRequired
	}

	tzName := strings.TrimSpace(query.Timezone)
	if tzName == "" {
		tzName = DefaultTimezone
	}

	loc, err := time.LoadLocation(tzName)
	if err != nil {
		loc = time.UTC
		tzName = "UTC"
	}

	day := query.Date
	localDay := normalizeDashboardDate(day, loc)
	startUTC := localDay.UTC()
	endUTC := localDay.Add(24*time.Hour - time.Nanosecond).UTC()

	defs, err := s.RecordRepository.ListMetricDefinitions(ctx, userID)
	if err != nil {
		return domain.DashboardSnapshot{}, err
	}

	records, err := s.RecordRepository.ListAllBetween(ctx, userID, startUTC, endUTC, DefaultDashboardLimit)
	if err != nil {
		return domain.DashboardSnapshot{}, err
	}

	goals, err := s.RecordRepository.ListGoalTemplates(ctx, userID)
	if err != nil {
		return domain.DashboardSnapshot{}, err
	}

	metrics := make([]domain.DashboardMetricValue, 0, len(defs))
	metricMap := make(map[string]domain.DashboardMetricValue, len(defs))
	for _, def := range defs {
		value := computeMetricValue(records, def)
		progress := 0.0
		if def.GoalDefault != nil && *def.GoalDefault > 0 {
			progress = clampPct((value / *def.GoalDefault) * 100)
		}

		status := DashboardMetricStatusPending
		if def.GoalDefault == nil {
			status = DashboardMetricStatusTracked
		} else if value >= *def.GoalDefault {
			status = DashboardMetricStatusCompleted
		}

		item := domain.DashboardMetricValue{
			MetricKey:   def.MetricKey,
			Label:       def.DisplayName,
			Value:       value,
			Unit:        def.Unit,
			Target:      def.GoalDefault,
			ProgressPct: progress,
			Status:      status,
		}
		metrics = append(metrics, item)
		metricMap[def.MetricKey] = item
	}

	goalValues := make([]domain.DashboardGoalValue, 0, len(goals))
	for _, goal := range goals {
		metric, ok := metricMap[goal.MetricKey]
		current := 0.0
		if ok {
			current = metric.Value
		}

		status, progress := evaluateGoal(current, goal.TargetValue, goal.Comparison)
		goalValues = append(goalValues, domain.DashboardGoalValue{
			GoalID:      goal.ID,
			Title:       goal.Title,
			MetricKey:   goal.MetricKey,
			Current:     current,
			Target:      goal.TargetValue,
			ProgressPct: progress,
			Status:      status,
		})
	}

	return domain.DashboardSnapshot{
		Date:     localDay,
		Timezone: tzName,
		Metrics:  metrics,
		Goals:    goalValues,
	}, nil
}

// UpsertMetricDefinition creates/updates a metric definition.
func (s *Service) UpsertMetricDefinition(ctx context.Context, userID uint64, cmd input.UpsertMetricDefinitionCommand) (domain.MetricDefinition, error) {
	if userID == 0 {
		return domain.MetricDefinition{}, ErrUserIDIsRequired
	}
	if cmd.TagID == 0 {
		return domain.MetricDefinition{}, ErrTagIDIsRequired
	}
	if strings.TrimSpace(cmd.MetricKey) == "" {
		return domain.MetricDefinition{}, errors.New(ErrDashboardMetricKeyRequired)
	}
	if strings.TrimSpace(cmd.DisplayName) == "" {
		return domain.MetricDefinition{}, errors.New(ErrDashboardDisplayNameRequired)
	}

	active := true
	if cmd.IsActive != nil {
		active = *cmd.IsActive
	}

	def := domain.MetricDefinition{
		UserID:      userID,
		MetricKey:   strings.TrimSpace(cmd.MetricKey),
		DisplayName: strings.TrimSpace(cmd.DisplayName),
		CategoryID:  cmd.CategoryID,
		TagID:       cmd.TagID,
		TagIDs:      normalizeMetricTagIDs(cmd.TagID, cmd.TagIDs),
		ValueSource: normalizeOrDefault(cmd.ValueSource, DashboardValueSourceCount),
		Aggregation: normalizeOrDefault(cmd.Aggregation, DashboardAggregationSum),
		Unit:        normalizeOrDefault(cmd.Unit, DashboardUnitCount),
		GoalDefault: cmd.GoalDefault,
		IsActive:    active,
	}
	if cmd.ID != nil {
		def.ID = *cmd.ID
	}

	return s.RecordRepository.UpsertMetricDefinition(ctx, def)
}

// UpsertGoalTemplate creates/updates a goal template.
func (s *Service) UpsertGoalTemplate(ctx context.Context, userID uint64, cmd input.UpsertGoalTemplateCommand) (domain.GoalTemplate, error) {
	if userID == 0 {
		return domain.GoalTemplate{}, ErrUserIDIsRequired
	}
	if strings.TrimSpace(cmd.MetricKey) == "" {
		return domain.GoalTemplate{}, errors.New(ErrDashboardMetricKeyRequired)
	}
	if strings.TrimSpace(cmd.Title) == "" {
		return domain.GoalTemplate{}, errors.New(ErrDashboardTitleRequired)
	}
	if cmd.TargetValue <= 0 {
		return domain.GoalTemplate{}, errors.New(ErrDashboardTargetValueRequired)
	}

	active := true
	if cmd.IsActive != nil {
		active = *cmd.IsActive
	}

	template := domain.GoalTemplate{
		UserID:      userID,
		MetricKey:   strings.TrimSpace(cmd.MetricKey),
		Title:       strings.TrimSpace(cmd.Title),
		TargetValue: cmd.TargetValue,
		Comparison:  normalizeOrDefault(cmd.Comparison, DashboardGoalComparisonGTE),
		Period:      normalizeOrDefault(cmd.Period, DashboardGoalPeriodDay),
		IsActive:    active,
	}
	if cmd.ID != nil {
		template.ID = *cmd.ID
	}

	return s.RecordRepository.UpsertGoalTemplate(ctx, template)
}

// DeleteGoalTemplate disables a goal template.
func (s *Service) DeleteGoalTemplate(ctx context.Context, userID uint64, goalTemplateID uint64) error {
	if userID == 0 {
		return ErrUserIDIsRequired
	}
	if goalTemplateID == 0 {
		return errors.New(ErrDashboardGoalTemplateIDRequired)
	}
	return s.RecordRepository.DeleteGoalTemplate(ctx, userID, goalTemplateID)
}

func normalizeOrDefault(value string, fallback string) string {
	v := strings.TrimSpace(strings.ToLower(value))
	if v == "" {
		return fallback
	}
	return v
}

func computeMetricValue(records []domain.Record, def domain.MetricDefinition) float64 {
	var matched int
	sum := 0.0
	latest := 0.0
	latestTime := time.Time{}
	tagSet := buildMetricTagSet(def)

	for _, rec := range records {
		if _, ok := tagSet[rec.TagID]; !ok {
			continue
		}

		matched++
		currentValue := extractRecordValue(rec, def.ValueSource)
		sum += currentValue

		if rec.EventTime.After(latestTime) {
			latestTime = rec.EventTime
			latest = currentValue
		}
	}

	switch def.Aggregation {
	case DashboardAggregationCount:
		return float64(matched)
	case DashboardAggregationAvg:
		if matched == 0 {
			return 0
		}
		return sum / float64(matched)
	case DashboardAggregationLatest:
		return latest
	default:
		return sum
	}
}

func buildMetricTagSet(def domain.MetricDefinition) map[uint64]struct{} {
	tagSet := make(map[uint64]struct{}, len(def.TagIDs)+1)
	if len(def.TagIDs) == 0 {
		tagSet[def.TagID] = struct{}{}
		return tagSet
	}
	for _, id := range def.TagIDs {
		tagSet[id] = struct{}{}
	}
	return tagSet
}

func normalizeMetricTagIDs(primary uint64, tagIDs []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(tagIDs)+1)
	out := make([]uint64, 0, len(tagIDs)+1)

	if primary != 0 {
		seen[primary] = struct{}{}
		out = append(out, primary)
	}

	for _, id := range tagIDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}

	if len(out) == 0 && primary != 0 {
		return []uint64{primary}
	}
	return out
}

func extractRecordValue(rec domain.Record, valueSource string) float64 {
	switch valueSource {
	case DashboardValueSourceDuration:
		if rec.DurationSecs != nil {
			return float64(*rec.DurationSecs)
		}
		return 0
	case DashboardValueSourceRaw, DashboardValueSourceLatestValue:
		if rec.Value != nil {
			return *rec.Value
		}
		return 0
	default:
		return 1
	}
}

func evaluateGoal(current float64, target float64, comparison string) (string, float64) {
	if target <= 0 {
		return DashboardMetricStatusInvalid, 0
	}

	switch comparison {
	case DashboardGoalComparisonLTE:
		progress := clampPct((target / math.Max(current, 0.000001)) * 100)
		if current <= target {
			return DashboardMetricStatusCompleted, 100
		}
		return DashboardMetricStatusPending, progress
	case DashboardGoalComparisonEQ:
		if math.Abs(current-target) < 0.0001 {
			return DashboardMetricStatusCompleted, 100
		}
		progress := clampPct(100 - math.Abs(current-target)/target*100)
		return DashboardMetricStatusPending, progress
	default:
		progress := clampPct((current / target) * 100)
		if current >= target {
			return DashboardMetricStatusCompleted, progress
		}
		return DashboardMetricStatusPending, progress
	}
}

func clampPct(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 999 {
		return 999
	}
	return math.Round(value*100) / 100
}
