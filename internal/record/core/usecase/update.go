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
	targetCategoryID := existing.CategoryID
	if cmd.CategoryID != nil {
		targetCategoryID = *cmd.CategoryID
	}
	targetTagID := existing.TagID
	if cmd.TagID != nil {
		targetTagID = *cmd.TagID
	}

	// Ensure tag-category alignment. If tag is zero/missing, or belongs to a different category,
	// auto-create a category-specific tag and use it.
	finalTagID := targetTagID
	if finalTagID == 0 {
		// create a new tag for this category with generated name
		newTag := tagdomain.Tag{
			Name:       fmt.Sprintf("auto-tag-%d", time.Now().UTC().UnixNano()),
			UserID:     userID,
			CategoryID: targetCategoryID,
		}
		created, err := s.TagRepository.Create(ctx, newTag)
		if err != nil {
			return domain.Record{}, fmt.Errorf("create auto tag: %w", err)
		}
		finalTagID = created.ID
	} else {
		// fetch tag and ensure category match
		t, err := s.TagRepository.GetByID(ctx, finalTagID, userID)
		if err != nil {
			return domain.Record{}, fmt.Errorf("lookup tag: %w", err)
		}
		if t.ID == 0 {
			// tag id invalid for user — create a new one in the category
			newTag := tagdomain.Tag{
				Name:       fmt.Sprintf("auto-tag-%d", time.Now().UTC().UnixNano()),
				UserID:     userID,
				CategoryID: targetCategoryID,
			}
			created, err := s.TagRepository.Create(ctx, newTag)
			if err != nil {
				return domain.Record{}, fmt.Errorf("create auto tag: %w", err)
			}
			finalTagID = created.ID
		} else if t.CategoryID != targetCategoryID {
			// create a category-specific twin (prefer same name; fallback to derived)
			candidate := tagdomain.Tag{
				Name:       t.Name,
				UserID:     userID,
				CategoryID: targetCategoryID,
			}
			created, err := s.TagRepository.Create(ctx, candidate)
			if err != nil {
				alt := tagdomain.Tag{
					Name:       fmt.Sprintf("%s@cat-%d", t.Name, targetCategoryID),
					UserID:     userID,
					CategoryID: targetCategoryID,
				}
				created, err = s.TagRepository.Create(ctx, alt)
				if err != nil {
					return domain.Record{}, fmt.Errorf("create category tag: %w", err)
				}
			}
			finalTagID = created.ID
		}
	}

	// Apply patch-like updates to the entity
	if cmd.Title != nil {
		existing.Title = *cmd.Title
	}
	if cmd.Description != nil {
		existing.Description = cmd.Description
	}
	existing.CategoryID = targetCategoryID
	existing.TagID = finalTagID
	if cmd.EventTime != nil {
		existing.EventTime = *cmd.EventTime
	}
	if cmd.RecordedAt != nil {
		existing.RecordedAt = cmd.RecordedAt
	}
	if cmd.DurationSecs != nil {
		existing.DurationSecs = cmd.DurationSecs
	}
	if cmd.Value != nil {
		existing.Value = cmd.Value
	}
	if cmd.Source != nil {
		existing.Source = cmd.Source
	}
	if cmd.Timezone != nil {
		existing.Timezone = cmd.Timezone
	}
	if cmd.Status != nil {
		existing.Status = cmd.Status
	}

	updated, err := s.RecordRepository.Update(ctx, existing)
	if err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}
