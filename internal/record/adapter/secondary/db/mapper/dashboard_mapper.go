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
