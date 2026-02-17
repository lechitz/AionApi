package repository

import (
	"context"

	dbport "github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

func (r *RecordRepository) ListMetricDefinitions(ctx context.Context, userID uint64) ([]domain.MetricDefinition, error) {
	var rows []model.MetricDefinition
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("metric_key ASC").
		Find(&rows).Error(); err != nil {
		return nil, err
	}

	out := make([]domain.MetricDefinition, len(rows))
	ids := make([]uint64, 0, len(rows))
	for i := range rows {
		out[i] = mapper.MetricDefinitionFromDB(rows[i])
		ids = append(ids, rows[i].ID)
	}
	if len(ids) == 0 {
		return out, nil
	}

	var bindings []model.MetricDefinitionTagBinding
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND metric_definition_id IN ?", userID, ids).
		Order("metric_definition_id ASC, tag_id ASC").
		Find(&bindings).Error(); err != nil {
		return nil, err
	}

	bindingsByMetric := make(map[uint64][]uint64, len(ids))
	for _, b := range bindings {
		bindingsByMetric[b.MetricDefinitionID] = append(bindingsByMetric[b.MetricDefinitionID], b.TagID)
	}

	for i := range out {
		seen := map[uint64]struct{}{out[i].TagID: {}}
		tagIDs := []uint64{out[i].TagID}
		for _, extraTagID := range bindingsByMetric[out[i].ID] {
			if _, ok := seen[extraTagID]; ok {
				continue
			}
			seen[extraTagID] = struct{}{}
			tagIDs = append(tagIDs, extraTagID)
		}
		out[i].TagIDs = tagIDs
	}
	return out, nil
}

func (r *RecordRepository) UpsertMetricDefinition(ctx context.Context, definition domain.MetricDefinition) (domain.MetricDefinition, error) {
	row := mapper.MetricDefinitionToDB(definition)
	err := r.db.WithContext(ctx).Transaction(func(tx dbport.DB) error {
		if row.ID != 0 {
			if err := tx.Model(&model.MetricDefinition{}).
				Where("id = ? AND user_id = ?", row.ID, row.UserID).
				Updates(map[string]interface{}{
					"metric_key":   row.MetricKey,
					"display_name": row.DisplayName,
					"category_id":  row.CategoryID,
					"tag_id":       row.TagID,
					"value_source": row.ValueSource,
					"aggregation":  row.Aggregation,
					"unit":         row.Unit,
					"goal_default": row.GoalDefault,
					"is_active":    row.IsActive,
				}).Error(); err != nil {
				return err
			}
			if err := tx.Where("id = ? AND user_id = ?", row.ID, row.UserID).First(&row).Error(); err != nil {
				return err
			}
		} else {
			if err := tx.Create(&row).Error(); err != nil {
				return err
			}
		}
		return r.syncMetricDefinitionTagBindings(ctx, tx, row.UserID, row.ID, definition.TagIDs)
	})
	if err != nil {
		return domain.MetricDefinition{}, err
	}
	return mapper.MetricDefinitionFromDB(row), nil
}

func (r *RecordRepository) syncMetricDefinitionTagBindings(ctx context.Context, tx dbport.DB, userID uint64, metricDefinitionID uint64, tagIDs []uint64) error {
	if err := tx.WithContext(ctx).
		Where("user_id = ? AND metric_definition_id = ?", userID, metricDefinitionID).
		Delete(&model.MetricDefinitionTagBinding{}).Error(); err != nil {
		return err
	}

	if len(tagIDs) == 0 {
		return nil
	}

	seen := make(map[uint64]struct{}, len(tagIDs))
	toInsert := make([]model.MetricDefinitionTagBinding, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		if tagID == 0 {
			continue
		}
		if _, ok := seen[tagID]; ok {
			continue
		}
		seen[tagID] = struct{}{}
		toInsert = append(toInsert, model.MetricDefinitionTagBinding{
			UserID:             userID,
			MetricDefinitionID: metricDefinitionID,
			TagID:              tagID,
		})
	}

	if len(toInsert) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Create(&toInsert).Error()
}

func (r *RecordRepository) ListGoalTemplates(ctx context.Context, userID uint64) ([]domain.GoalTemplate, error) {
	var rows []model.GoalTemplate
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("id ASC").
		Find(&rows).Error(); err != nil {
		return nil, err
	}

	out := make([]domain.GoalTemplate, len(rows))
	for i := range rows {
		out[i] = mapper.GoalTemplateFromDB(rows[i])
	}
	return out, nil
}

func (r *RecordRepository) UpsertGoalTemplate(ctx context.Context, template domain.GoalTemplate) (domain.GoalTemplate, error) {
	row := mapper.GoalTemplateToDB(template)
	if row.ID != 0 {
		if err := r.db.WithContext(ctx).
			Model(&model.GoalTemplate{}).
			Where("id = ? AND user_id = ?", row.ID, row.UserID).
			Updates(map[string]interface{}{
				"metric_key":   row.MetricKey,
				"title":        row.Title,
				"target_value": row.TargetValue,
				"comparison":   row.Comparison,
				"period":       row.Period,
				"is_active":    row.IsActive,
			}).Error(); err != nil {
			return domain.GoalTemplate{}, err
		}
		if err := r.db.WithContext(ctx).
			Where("id = ? AND user_id = ?", row.ID, row.UserID).
			First(&row).Error(); err != nil {
			return domain.GoalTemplate{}, err
		}
		return mapper.GoalTemplateFromDB(row), nil
	}

	if err := r.db.WithContext(ctx).Create(&row).Error(); err != nil {
		return domain.GoalTemplate{}, err
	}
	return mapper.GoalTemplateFromDB(row), nil
}

func (r *RecordRepository) DeleteGoalTemplate(ctx context.Context, userID uint64, goalTemplateID uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.GoalTemplate{}).
		Where("id = ? AND user_id = ?", goalTemplateID, userID).
		Update("is_active", false).Error()
}
