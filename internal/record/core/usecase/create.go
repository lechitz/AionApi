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

// Create creates a new record after validating inputs.
func (s *Service) Create(ctx context.Context, cmd input.CreateRecordCommand) (domain.Record, error) {
	// title required
	if cmd.Title == "" {
		return domain.Record{}, errors.New("title required")
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return domain.Record{}, err
	}

	// default eventTime to recordedAt or now if zero
	eventTime := cmd.EventTime
	if eventTime.IsZero() {
		if cmd.RecordedAt != nil {
			eventTime = *cmd.RecordedAt
		} else {
			eventTime = time.Now().UTC()
		}
	}

	// validate recordedAt not future
	if cmd.RecordedAt != nil && cmd.RecordedAt.After(time.Now().UTC()) {
		return domain.Record{}, errors.New("recordedAt cannot be in the future")
	}

	// Tag validation/auto-creation logic
	finalTagID := cmd.TagID
	if finalTagID != 0 {
		// try fetch tag by id for this user
		tagObj, err := s.TagRepository.GetByID(ctx, finalTagID, userID)
		if err != nil {
			return domain.Record{}, fmt.Errorf("lookup tag: %w", err)
		}
		if tagObj.ID == 0 {
			// tag not found — create one for the specified category
			newTag := tagdomain.Tag{
				Name:       fmt.Sprintf("tag-%d", finalTagID),
				UserID:     userID,
				CategoryID: cmd.CategoryID,
			}
			createdTag, err := s.TagRepository.Create(ctx, newTag)
			if err != nil {
				return domain.Record{}, fmt.Errorf("create fallback tag: %w", err)
			}
			finalTagID = createdTag.ID
		} else {
			// tag exists — ensure it belongs to the category
			if tagObj.CategoryID != cmd.CategoryID {
				return domain.Record{}, errors.New("tag belongs to a different category")
			}
		}
	} else {
		// TagID is zero — create a new tag for this category with a generated name
		newTag := tagdomain.Tag{
			Name:       fmt.Sprintf("auto-tag-%d", time.Now().UTC().UnixNano()),
			UserID:     userID,
			CategoryID: cmd.CategoryID,
		}
		createdTag, err := s.TagRepository.Create(ctx, newTag)
		if err != nil {
			return domain.Record{}, fmt.Errorf("create auto tag: %w", err)
		}
		finalTagID = createdTag.ID
	}

	// construct domain.Record
	rec := domain.Record{
		UserID:       userID,
		Title:        cmd.Title,
		Description:  cmd.Description,
		CategoryID:   cmd.CategoryID,
		EventTime:    eventTime,
		RecordedAt:   cmd.RecordedAt,
		DurationSecs: cmd.DurationSecs,
		Value:        cmd.Value,
		Source:       cmd.Source,
		Timezone:     cmd.Timezone,
		Status:       cmd.Status,
		TagID:        finalTagID,
	}

	created, err := s.RecordRepository.Create(ctx, rec)
	if err != nil {
		return domain.Record{}, err
	}
	return created, nil
}
