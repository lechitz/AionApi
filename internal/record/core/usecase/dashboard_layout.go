package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
)

func (s *Service) ListDashboardViews(ctx context.Context, userID uint64) ([]domain.DashboardView, error) {
	if userID == 0 {
		return nil, ErrUserIDIsRequired
	}
	return s.ensureDashboardViews(ctx, userID)
}

func (s *Service) GetDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error) {
	if userID == 0 {
		return domain.DashboardView{}, ErrUserIDIsRequired
	}
	if viewID == 0 {
		return domain.DashboardView{}, errors.New("viewID is required")
	}
	return s.RecordRepository.GetDashboardView(ctx, userID, viewID)
}

func (s *Service) CreateDashboardView(ctx context.Context, userID uint64, cmd input.CreateDashboardViewCommand) (domain.DashboardView, error) {
	if userID == 0 {
		return domain.DashboardView{}, ErrUserIDIsRequired
	}

	name := strings.TrimSpace(cmd.Name)
	if name == "" {
		name = "Meu Dashboard"
	}
	isDefault := cmd.IsDefault != nil && *cmd.IsDefault

	view, err := s.RecordRepository.CreateDashboardView(ctx, domain.DashboardView{
		UserID:    userID,
		Name:      name,
		IsDefault: isDefault,
	})
	if err != nil {
		return domain.DashboardView{}, err
	}

	if isDefault {
		return s.RecordRepository.SetDefaultDashboardView(ctx, userID, view.ID)
	}
	return view, nil
}

func (s *Service) SetDefaultDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error) {
	if userID == 0 {
		return domain.DashboardView{}, ErrUserIDIsRequired
	}
	if viewID == 0 {
		return domain.DashboardView{}, errors.New("viewID is required")
	}
	return s.RecordRepository.SetDefaultDashboardView(ctx, userID, viewID)
}

func (s *Service) UpsertDashboardWidget(ctx context.Context, userID uint64, cmd input.UpsertDashboardWidgetCommand) (domain.DashboardWidget, error) {
	if userID == 0 {
		return domain.DashboardWidget{}, ErrUserIDIsRequired
	}
	if cmd.ViewID == 0 {
		return domain.DashboardWidget{}, errors.New("viewID is required")
	}
	if cmd.MetricDefinitionID == 0 {
		return domain.DashboardWidget{}, errors.New("metricDefinitionID is required")
	}

	widgetType := normalizeWidgetType(cmd.WidgetType)
	size := normalizeWidgetSize(cmd.Size)

	if size == domain.DashboardWidgetSizeLarge {
		var exclude *uint64
		if cmd.ID != nil && *cmd.ID != 0 {
			exclude = cmd.ID
		}
		countLarge, err := s.RecordRepository.CountLargeWidgetsInView(ctx, userID, cmd.ViewID, exclude)
		if err != nil {
			return domain.DashboardWidget{}, err
		}
		if countLarge >= domain.MaxLargeWidgetsPerDashboard {
			return domain.DashboardWidget{}, fmt.Errorf("limit reached: max %d large widgets", domain.MaxLargeWidgetsPerDashboard)
		}
	}

	orderIndex := 0
	if cmd.OrderIndex != nil {
		orderIndex = *cmd.OrderIndex
	} else {
		existing, err := s.RecordRepository.ListDashboardWidgetsByView(ctx, userID, cmd.ViewID)
		if err != nil {
			return domain.DashboardWidget{}, err
		}
		orderIndex = len(existing)
	}

	isActive := true
	if cmd.IsActive != nil {
		isActive = *cmd.IsActive
	}

	widget := domain.DashboardWidget{
		UserID:             userID,
		ViewID:             cmd.ViewID,
		MetricDefinitionID: cmd.MetricDefinitionID,
		WidgetType:         widgetType,
		Size:               size,
		OrderIndex:         orderIndex,
		TitleOverride:      cmd.TitleOverride,
		ConfigJSON:         strings.TrimSpace(cmd.ConfigJSON),
		IsActive:           isActive,
	}
	if widget.ConfigJSON == "" {
		widget.ConfigJSON = "{}"
	}
	if cmd.ID != nil {
		widget.ID = *cmd.ID
	}

	return s.RecordRepository.UpsertDashboardWidget(ctx, widget)
}

func (s *Service) ReorderDashboardWidgets(ctx context.Context, userID uint64, cmd input.ReorderDashboardWidgetsCommand) ([]domain.DashboardWidget, error) {
	if userID == 0 {
		return nil, ErrUserIDIsRequired
	}
	if cmd.ViewID == 0 {
		return nil, errors.New("viewID is required")
	}
	if len(cmd.Items) == 0 {
		return nil, errors.New("items are required")
	}

	items := make([]domain.DashboardWidget, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		if item.WidgetID == 0 {
			return nil, errors.New("widgetID is required")
		}
		items = append(items, domain.DashboardWidget{
			ID:         item.WidgetID,
			OrderIndex: item.OrderIndex,
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].OrderIndex < items[j].OrderIndex })
	return s.RecordRepository.ReorderDashboardWidgets(ctx, userID, cmd.ViewID, items)
}

func (s *Service) DeleteDashboardWidget(ctx context.Context, userID uint64, widgetID uint64) error {
	if userID == 0 {
		return ErrUserIDIsRequired
	}
	if widgetID == 0 {
		return errors.New("widgetID is required")
	}
	return s.RecordRepository.DeleteDashboardWidget(ctx, userID, widgetID)
}

func (s *Service) CreateMetricAndWidget(ctx context.Context, userID uint64, cmd input.CreateMetricAndWidgetCommand) (domain.DashboardWidget, error) {
	metric, err := s.UpsertMetricDefinition(ctx, userID, cmd.Metric)
	if err != nil {
		return domain.DashboardWidget{}, err
	}

	widgetCmd := cmd.Widget
	if widgetCmd.MetricDefinitionID == 0 {
		widgetCmd.MetricDefinitionID = metric.ID
	}
	return s.UpsertDashboardWidget(ctx, userID, widgetCmd)
}

func (s *Service) SuggestMetricDefinitions(ctx context.Context, userID uint64, limit int) ([]domain.MetricDefinitionSuggestion, error) {
	if userID == 0 {
		return nil, ErrUserIDIsRequired
	}
	if limit <= 0 {
		limit = 8
	}
	if limit > 20 {
		limit = 20
	}

	// Deterministic suggestions from existing active tags.
	tags, err := s.TagRepository.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return []domain.MetricDefinitionSuggestion{}, nil
	}

	out := make([]domain.MetricDefinitionSuggestion, 0, limit)
	seen := make(map[string]struct{}, limit)
	for _, tag := range tags {
		if len(out) >= limit {
			break
		}
		key := slugMetricKey(tag.Name)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		tagID := tag.ID
		out = append(out, domain.MetricDefinitionSuggestion{
			MetricKey:   key,
			DisplayName: strings.TrimSpace(tag.Name),
			CategoryID:  &tag.CategoryID,
			TagIDs:      []uint64{tagID},
			ValueSource: "count",
			Aggregation: "sum",
			Unit:        "count",
			Reason:      "Baseado em tags existentes da sua taxonomia.",
		})
	}
	return out, nil
}

func (s *Service) ensureDashboardViews(ctx context.Context, userID uint64) ([]domain.DashboardView, error) {
	views, err := s.RecordRepository.ListDashboardViews(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(views) > 0 {
		return views, nil
	}

	defaultView, err := s.RecordRepository.CreateDashboardView(ctx, domain.DashboardView{
		UserID:    userID,
		Name:      "Principal",
		IsDefault: true,
	})
	if err != nil {
		return nil, err
	}
	_ = defaultView

	views, err = s.RecordRepository.ListDashboardViews(ctx, userID)
	if err != nil {
		return nil, err
	}
	return views, nil
}

func normalizeWidgetType(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case domain.DashboardWidgetTypeGoalProgress:
		return domain.DashboardWidgetTypeGoalProgress
	case domain.DashboardWidgetTypeTrendLine:
		return domain.DashboardWidgetTypeTrendLine
	case domain.DashboardWidgetTypeChecklist:
		return domain.DashboardWidgetTypeChecklist
	default:
		return domain.DashboardWidgetTypeKPINumber
	}
}

func normalizeWidgetSize(v string) string {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case domain.DashboardWidgetSizeMedium:
		return domain.DashboardWidgetSizeMedium
	case domain.DashboardWidgetSizeLarge:
		return domain.DashboardWidgetSizeLarge
	default:
		return domain.DashboardWidgetSizeSmall
	}
}

func slugMetricKey(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if v == "" {
		return ""
	}
	replacer := strings.NewReplacer(
		" ", "_",
		"-", "_",
		"/", "_",
		"(", "",
		")", "",
	)
	v = replacer.Replace(v)
	return strings.Trim(v, "_")
}
