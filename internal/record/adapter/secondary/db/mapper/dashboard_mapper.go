package mapper

import (
	dbmodel "github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

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

func DashboardViewToDB(in domain.DashboardView) dbmodel.DashboardView {
	return dbmodel.DashboardView{
		ID:        in.ID,
		UserID:    in.UserID,
		Name:      in.Name,
		IsDefault: in.IsDefault,
	}
}

func DashboardWidgetFromDB(in dbmodel.DashboardWidget) domain.DashboardWidget {
	return domain.DashboardWidget{
		ID:                 in.ID,
		UserID:             in.UserID,
		ViewID:             in.ViewID,
		MetricDefinitionID: in.MetricDefinitionID,
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

func DashboardWidgetToDB(in domain.DashboardWidget) dbmodel.DashboardWidget {
	return dbmodel.DashboardWidget{
		ID:                 in.ID,
		UserID:             in.UserID,
		ViewID:             in.ViewID,
		MetricDefinitionID: in.MetricDefinitionID,
		WidgetType:         in.WidgetType,
		Size:               in.Size,
		OrderIndex:         in.OrderIndex,
		TitleOverride:      in.TitleOverride,
		ConfigJSON:         in.ConfigJSON,
		IsActive:           in.IsActive,
	}
}
