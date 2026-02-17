package repository

import (
	"context"
	"fmt"

	dbport "github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

func (r *RecordRepository) ListDashboardViews(ctx context.Context, userID uint64) ([]domain.DashboardView, error) {
	var rows []model.DashboardView
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, id ASC").
		Find(&rows).Error(); err != nil {
		return nil, err
	}

	out := make([]domain.DashboardView, len(rows))
	for i := range rows {
		out[i] = mapper.DashboardViewFromDB(rows[i])
	}
	return out, nil
}

func (r *RecordRepository) GetDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error) {
	var row model.DashboardView
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", viewID, userID).
		First(&row).Error(); err != nil {
		return domain.DashboardView{}, err
	}

	view := mapper.DashboardViewFromDB(row)
	widgets, err := r.ListDashboardWidgetsByView(ctx, userID, viewID)
	if err != nil {
		return domain.DashboardView{}, err
	}
	view.Widgets = widgets
	return view, nil
}

func (r *RecordRepository) CreateDashboardView(ctx context.Context, view domain.DashboardView) (domain.DashboardView, error) {
	row := mapper.DashboardViewToDB(view)
	if err := r.db.WithContext(ctx).Create(&row).Error(); err != nil {
		return domain.DashboardView{}, err
	}
	return mapper.DashboardViewFromDB(row), nil
}

func (r *RecordRepository) SetDefaultDashboardView(ctx context.Context, userID uint64, viewID uint64) (domain.DashboardView, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx dbport.DB) error {
		if err := tx.Model(&model.DashboardView{}).
			Where("user_id = ?", userID).
			Update("is_default", false).Error(); err != nil {
			return err
		}

		if err := tx.Model(&model.DashboardView{}).
			Where("id = ? AND user_id = ?", viewID, userID).
			Update("is_default", true).Error(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return domain.DashboardView{}, err
	}

	var row model.DashboardView
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", viewID, userID).
		First(&row).Error(); err != nil {
		return domain.DashboardView{}, err
	}
	return mapper.DashboardViewFromDB(row), nil
}

func (r *RecordRepository) UpsertDashboardWidget(ctx context.Context, widget domain.DashboardWidget) (domain.DashboardWidget, error) {
	row := mapper.DashboardWidgetToDB(widget)
	if row.ConfigJSON == "" {
		row.ConfigJSON = "{}"
	}

	if row.ID != 0 {
		if err := r.db.WithContext(ctx).
			Model(&model.DashboardWidget{}).
			Where("id = ? AND user_id = ?", row.ID, row.UserID).
			Updates(map[string]interface{}{
				"view_id":              row.ViewID,
				"metric_definition_id": row.MetricDefinitionID,
				"widget_type":          row.WidgetType,
				"size":                 row.Size,
				"order_index":          row.OrderIndex,
				"title_override":       row.TitleOverride,
				"config_json":          row.ConfigJSON,
				"is_active":            row.IsActive,
			}).Error(); err != nil {
			return domain.DashboardWidget{}, err
		}
		if err := r.db.WithContext(ctx).
			Where("id = ? AND user_id = ?", row.ID, row.UserID).
			First(&row).Error(); err != nil {
			return domain.DashboardWidget{}, err
		}
		return mapper.DashboardWidgetFromDB(row), nil
	}

	if err := r.db.WithContext(ctx).Create(&row).Error(); err != nil {
		return domain.DashboardWidget{}, err
	}
	return mapper.DashboardWidgetFromDB(row), nil
}

func (r *RecordRepository) ListDashboardWidgetsByView(ctx context.Context, userID uint64, viewID uint64) ([]domain.DashboardWidget, error) {
	var rows []model.DashboardWidget
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND view_id = ? AND is_active = ?", userID, viewID, true).
		Order("order_index ASC, id ASC").
		Find(&rows).Error(); err != nil {
		return nil, err
	}

	out := make([]domain.DashboardWidget, len(rows))
	for i := range rows {
		out[i] = mapper.DashboardWidgetFromDB(rows[i])
	}
	return out, nil
}

func (r *RecordRepository) ReorderDashboardWidgets(ctx context.Context, userID uint64, viewID uint64, items []domain.DashboardWidget) ([]domain.DashboardWidget, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx dbport.DB) error {
		for _, item := range items {
			if err := tx.Model(&model.DashboardWidget{}).
				Where("id = ? AND user_id = ? AND view_id = ?", item.ID, userID, viewID).
				Update("order_index", item.OrderIndex).Error(); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.ListDashboardWidgetsByView(ctx, userID, viewID)
}

func (r *RecordRepository) DeleteDashboardWidget(ctx context.Context, userID uint64, widgetID uint64) error {
	return r.db.WithContext(ctx).
		Model(&model.DashboardWidget{}).
		Where("id = ? AND user_id = ?", widgetID, userID).
		Update("is_active", false).Error()
}

func (r *RecordRepository) CountLargeWidgetsInView(ctx context.Context, userID uint64, viewID uint64, excludeWidgetID *uint64) (int64, error) {
	q := r.db.WithContext(ctx).
		Model(&model.DashboardWidget{}).
		Where("user_id = ? AND view_id = ? AND is_active = ? AND size = ?", userID, viewID, true, domain.DashboardWidgetSizeLarge)
	if excludeWidgetID != nil && *excludeWidgetID != 0 {
		q = q.Where("id <> ?", *excludeWidgetID)
	}

	var count int64
	if err := q.Count(&count).Error(); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RecordRepository) nextWidgetOrderIndex(ctx context.Context, userID uint64, viewID uint64) (int, error) {
	type out struct {
		MaxOrder *int `gorm:"column:max_order"`
	}
	var row out
	if err := r.db.WithContext(ctx).
		Model(&model.DashboardWidget{}).
		Select("MAX(order_index) AS max_order").
		Where("user_id = ? AND view_id = ? AND is_active = ?", userID, viewID, true).
		Scan(&row).Error(); err != nil {
		return 0, err
	}
	if row.MaxOrder == nil {
		return 0, nil
	}
	return *row.MaxOrder + 1, nil
}

func (r *RecordRepository) widgetBelongsToUserView(ctx context.Context, userID uint64, viewID uint64, widgetID uint64) error {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.DashboardWidget{}).
		Where("id = ? AND user_id = ? AND view_id = ?", widgetID, userID, viewID).
		Count(&count).Error(); err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("widget %d not found for user/view", widgetID)
	}
	return nil
}
