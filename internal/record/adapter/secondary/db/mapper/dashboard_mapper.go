package mapper

import (
	dbmodel "github.com/lechitz/aion-api/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/record/core/domain"
)

// MetricDefinitionFromDB maps a DB metric definition row into the core domain model.
func MetricDefinitionFromDB(in dbmodel.MetricDefinition) domain.MetricDefinition {
	return domain.MetricDefinition{
		ID:          in.ID,
		UserID:      in.UserID,
		MetricKey:   in.MetricKey,
		DisplayName: in.DisplayName,
		CategoryID:  in.CategoryID,
		TagID:       in.TagID,
		TagIDs:      []uint64{in.TagID},
		ValueSource: in.ValueSource,
		Aggregation: in.Aggregation,
		Unit:        in.Unit,
		GoalDefault: in.GoalDefault,
		IsActive:    in.IsActive,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}

// MetricDefinitionToDB maps a core metric definition into the DB persistence model.
func MetricDefinitionToDB(in domain.MetricDefinition) dbmodel.MetricDefinition {
	return dbmodel.MetricDefinition{
		ID:          in.ID,
		UserID:      in.UserID,
		MetricKey:   in.MetricKey,
		DisplayName: in.DisplayName,
		CategoryID:  in.CategoryID,
		TagID:       in.TagID,
		ValueSource: in.ValueSource,
		Aggregation: in.Aggregation,
		Unit:        in.Unit,
		GoalDefault: in.GoalDefault,
		IsActive:    in.IsActive,
	}
}

// GoalTemplateFromDB maps a DB goal template row into the core domain model.
func GoalTemplateFromDB(in dbmodel.GoalTemplate) domain.GoalTemplate {
	return domain.GoalTemplate{
		ID:          in.ID,
		UserID:      in.UserID,
		MetricKey:   in.MetricKey,
		Title:       in.Title,
		TargetValue: in.TargetValue,
		Comparison:  in.Comparison,
		Period:      in.Period,
		IsActive:    in.IsActive,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}

// GoalTemplateToDB maps a core goal template into the DB persistence model.
func GoalTemplateToDB(in domain.GoalTemplate) dbmodel.GoalTemplate {
	return dbmodel.GoalTemplate{
		ID:          in.ID,
		UserID:      in.UserID,
		MetricKey:   in.MetricKey,
		Title:       in.Title,
		TargetValue: in.TargetValue,
		Comparison:  in.Comparison,
		Period:      in.Period,
		IsActive:    in.IsActive,
	}
}

// DashboardViewFromDB maps a DB dashboard view row into the core domain model.
func DashboardViewFromDB(in dbmodel.DashboardView) domain.DashboardView {
	return domain.DashboardView{
		ID:        in.ID,
		UserID:    in.UserID,
		Name:      in.Name,
		IsDefault: in.IsDefault,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

// DashboardViewToDB maps a core dashboard view into the DB persistence model.
func DashboardViewToDB(in domain.DashboardView) dbmodel.DashboardView {
	return dbmodel.DashboardView{
		ID:        in.ID,
		UserID:    in.UserID,
		Name:      in.Name,
		IsDefault: in.IsDefault,
	}
}

// DashboardWidgetFromDB maps a DB dashboard widget row into the core domain model.
func DashboardWidgetFromDB(in dbmodel.DashboardWidget) domain.DashboardWidget {
	var metricDefinitionID uint64
	if in.MetricDefinitionID != nil {
		metricDefinitionID = *in.MetricDefinitionID
	}
	return domain.DashboardWidget{
		ID:                 in.ID,
		UserID:             in.UserID,
		ViewID:             in.ViewID,
		MetricDefinitionID: metricDefinitionID,
		WidgetType:         in.WidgetType,
		Size:               in.Size,
		OrderIndex:         in.OrderIndex,
		TitleOverride:      in.TitleOverride,
		ConfigJSON:         in.ConfigJSON,
		IsActive:           in.IsActive,
		CreatedAt:          in.CreatedAt,
		UpdatedAt:          in.UpdatedAt,
	}
}

// DashboardWidgetToDB maps a core dashboard widget into the DB persistence model.
func DashboardWidgetToDB(in domain.DashboardWidget) dbmodel.DashboardWidget {
	var metricDefinitionID *uint64
	if in.MetricDefinitionID != 0 {
		value := in.MetricDefinitionID
		metricDefinitionID = &value
	}
	return dbmodel.DashboardWidget{
		ID:                 in.ID,
		UserID:             in.UserID,
		ViewID:             in.ViewID,
		MetricDefinitionID: metricDefinitionID,
		WidgetType:         in.WidgetType,
		Size:               in.Size,
		OrderIndex:         in.OrderIndex,
		TitleOverride:      in.TitleOverride,
		ConfigJSON:         in.ConfigJSON,
		IsActive:           in.IsActive,
	}
}
