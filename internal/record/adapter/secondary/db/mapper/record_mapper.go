// Package mapper provides utility functions for mapping between domain and database objects.
package mapper

import (
	"github.com/lechitz/AionApi/internal/record/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/record/core/domain"
)

// RecordFromDB maps a model.Record object to a domain.Record object.
func RecordFromDB(record model.Record) domain.Record {
	return domain.Record{
		ID:           record.ID,
		UserID:       record.UserID,
		Title:        record.Title,
		Description:  record.Description,
		CategoryID:   record.CategoryID,
		TagID:        record.TagID,
		EventTime:    record.EventTime,
		RecordedAt:   record.RecordedAt,
		DurationSecs: record.DurationSecs,
		Value:        record.Value,
		Source:       record.Source,
		Timezone:     record.Timezone,
		Status:       record.Status,
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
		DeletedAt:    record.DeletedAt,
	}
}

// RecordToDB maps a domain.Record object to a model.Record object for database operations.
func RecordToDB(record domain.Record) model.Record {
	return model.Record{
		ID:           record.ID,
		UserID:       record.UserID,
		Title:        record.Title,
		Description:  record.Description,
		CategoryID:   record.CategoryID,
		TagID:        record.TagID,
		EventTime:    record.EventTime,
		RecordedAt:   record.RecordedAt,
		DurationSecs: record.DurationSecs,
		Value:        record.Value,
		Source:       record.Source,
		Timezone:     record.Timezone,
		Status:       record.Status,
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
		DeletedAt:    record.DeletedAt,
	}
}

// RecordsFromDB converts a slice of model.Record to a slice of domain.Record.
func RecordsFromDB(records []model.Record) []domain.Record {
	result := make([]domain.Record, len(records))
	for i, rec := range records {
		result[i] = RecordFromDB(rec)
	}
	return result
}
