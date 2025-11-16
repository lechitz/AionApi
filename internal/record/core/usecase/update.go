package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
	tagdomain "github.com/lechitz/AionApi/internal/tag/core/domain"
)

// Update applies partial changes to an existing record owned by the user.
func (s *Service) Update(ctx context.Context, recordID uint64, userID uint64, cmd input.UpdateRecordCommand) (domain.Record, error) {
	if recordID == 0 || userID == 0 {
		return domain.Record{}, errors.New("invalid recordID or userID")
	}

	existing, err := s.RecordRepository.GetByID(ctx, recordID, userID)
	if err != nil {
		return domain.Record{}, err
	}

	// Determine target category and tag after patch
	targetCategoryID := resolveUint64(existing.CategoryID, cmd.CategoryID)
	targetTagID := resolveUint64(existing.TagID, cmd.TagID)

	// Ensure tag-category alignment and pick the final tag id
	finalTagID, err := s.ensureTagCategoryAlignment(ctx, userID, targetCategoryID, targetTagID)
	if err != nil {
		return domain.Record{}, err
	}

	// Apply patch-like updates to the entity
	existing = applyRecordPatch(existing, cmd, targetCategoryID, finalTagID)

	updated, err := s.RecordRepository.Update(ctx, existing)
	if err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}

// resolveUint64 returns ptr value if provided, otherwise the fallback.
func resolveUint64(fallback uint64, ptr *uint64) uint64 {
	if ptr != nil {
		return *ptr
	}
	return fallback
}

// ensureTagCategoryAlignment guarantees we have a tag ID that belongs to the given category.
// If tagID is zero/invalid or belongs to a different category, it creates an appropriate tag.
func (s *Service) ensureTagCategoryAlignment(ctx context.Context, userID, categoryID, tagID uint64) (uint64, error) {
	if tagID == 0 {
		created, err := s.TagRepository.Create(ctx, tagdomain.Tag{
			Name:       fmt.Sprintf("auto-tag-%d", time.Now().UTC().UnixNano()),
			UserID:     userID,
			CategoryID: categoryID,
		})
		if err != nil {
			return 0, fmt.Errorf("create auto tag: %w", err)
		}
		return created.ID, nil
	}

	// fetch tag and ensure category match
	t, err := s.TagRepository.GetByID(ctx, tagID, userID)
	if err != nil {
		return 0, fmt.Errorf("lookup tag: %w", err)
	}
	if t.ID == 0 {
		created, err := s.TagRepository.Create(ctx, tagdomain.Tag{
			Name:       fmt.Sprintf("auto-tag-%d", time.Now().UTC().UnixNano()),
			UserID:     userID,
			CategoryID: categoryID,
		})
		if err != nil {
			return 0, fmt.Errorf("create auto tag: %w", err)
		}
		return created.ID, nil
	}
	if t.CategoryID != categoryID {
		candidate := tagdomain.Tag{
			Name:       t.Name,
			UserID:     userID,
			CategoryID: categoryID,
		}
		created, err := s.TagRepository.Create(ctx, candidate)
		if err != nil {
			alt := tagdomain.Tag{
				Name:       fmt.Sprintf("%s@cat-%d", t.Name, categoryID),
				UserID:     userID,
				CategoryID: categoryID,
			}
			created, err = s.TagRepository.Create(ctx, alt)
			if err != nil {
				return 0, fmt.Errorf("create category tag: %w", err)
			}
		}
		return created.ID, nil
	}
	return tagID, nil
}

// applyRecordPatch mutates a copy of the record with fields from cmd and the resolved IDs.
func applyRecordPatch(r domain.Record, cmd input.UpdateRecordCommand, categoryID, tagID uint64) domain.Record {
	if cmd.Title != nil {
		r.Title = *cmd.Title
	}
	if cmd.Description != nil {
		r.Description = cmd.Description
	}
	r.CategoryID = categoryID
	r.TagID = tagID
	if cmd.EventTime != nil {
		r.EventTime = *cmd.EventTime
	}
	if cmd.RecordedAt != nil {
		r.RecordedAt = cmd.RecordedAt
	}
	if cmd.DurationSecs != nil {
		r.DurationSecs = cmd.DurationSecs
	}
	if cmd.Value != nil {
		r.Value = cmd.Value
	}
	if cmd.Source != nil {
		r.Source = cmd.Source
	}
	if cmd.Timezone != nil {
		r.Timezone = cmd.Timezone
	}
	if cmd.Status != nil {
		r.Status = cmd.Status
	}
	return r
}
