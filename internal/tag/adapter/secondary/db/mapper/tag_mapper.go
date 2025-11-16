// Package mapper provides utility functions for mapping between domain and database objects.
package mapper

import (
	"time"

	"github.com/lechitz/AionApi/internal/tag/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/tag/core/domain"
)

// TagFromDB maps a TagDB to a domain.Tag.
func TagFromDB(db model.TagDB) domain.Tag {
	var deletedAt *time.Time
	if db.DeletedAt.Valid {
		deletedAt = &db.DeletedAt.Time
	}
	return domain.Tag{
		ID:          db.ID,
		UserID:      db.UserID,
		CategoryID:  db.CategoryID,
		Name:        db.Name,
		Description: db.Description,
		CreatedAt:   db.CreatedAt,
		UpdatedAt:   db.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

// TagsFromDB maps a slice of TagDB to a slice of domain.Tag.
func TagsFromDB(dbTags []model.TagDB) []domain.Tag {
	tags := make([]domain.Tag, len(dbTags))
	for i, db := range dbTags {
		tags[i] = TagFromDB(db)
	}
	return tags
}

// TagToDB maps a domain.Tag to TagDB.
func TagToDB(t domain.Tag) model.TagDB {
	return model.TagDB{
		ID:          t.ID,
		UserID:      t.UserID,
		CategoryID:  t.CategoryID,
		Name:        t.Name,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
