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
	if cmd.Title == "" {
		return domain.Record{}, errors.New("title required")
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return domain.Record{}, err
	}

	eventTime := resolveEventTime(cmd)

	if err := validateRecordedAt(cmd.RecordedAt); err != nil {
		return domain.Record{}, err
	}

	finalTagID, err := s.resolveTagID(ctx, cmd.TagID, cmd.CategoryID, userID)
	if err != nil {
		return domain.Record{}, err
	}

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

// resolveEventTime determines the event time from the command.
func resolveEventTime(cmd input.CreateRecordCommand) time.Time {
	if !cmd.EventTime.IsZero() {
		return cmd.EventTime
	}
	if cmd.RecordedAt != nil {
		return *cmd.RecordedAt
	}
	return time.Now().UTC()
}

// validateRecordedAt ensures recordedAt is not in the future.
func validateRecordedAt(recordedAt *time.Time) error {
	if recordedAt != nil && recordedAt.After(time.Now().UTC()) {
		return errors.New("recordedAt cannot be in the future")
	}
	return nil
}

// resolveTagID handles tag validation and auto-creation logic.
func (s *Service) resolveTagID(ctx context.Context, tagID, categoryID uint64, userID uint64) (uint64, error) {
	if tagID != 0 {
		return s.resolveExistingTag(ctx, tagID, categoryID, userID)
	}
	return s.createAutoTag(ctx, categoryID, userID)
}

// resolveExistingTag validates an existing tag or creates a category-specific one.
func (s *Service) resolveExistingTag(ctx context.Context, tagID, categoryID uint64, userID uint64) (uint64, error) {
	tagObj, err := s.TagRepository.GetByID(ctx, tagID, userID)
	if err != nil {
		return 0, fmt.Errorf("lookup tag: %w", err)
	}

	if tagObj.ID == 0 {
		return s.createAutoTag(ctx, categoryID, userID)
	}

	if tagObj.CategoryID != categoryID {
		return s.createCategorySpecificTag(ctx, tagObj, categoryID, userID)
	}

	return tagID, nil
}

// createCategorySpecificTag creates a tag for a specific category.
func (s *Service) createCategorySpecificTag(ctx context.Context, original tagdomain.Tag, categoryID uint64, userID uint64) (uint64, error) {
	candidate := tagdomain.Tag{
		Name:       original.Name,
		UserID:     userID,
		CategoryID: categoryID,
	}

	createdTag, err := s.TagRepository.Create(ctx, candidate)
	if err != nil {
		// Fallback to a derived unique name per category
		alt := tagdomain.Tag{
			Name:       fmt.Sprintf("%s@cat-%d", original.Name, categoryID),
			UserID:     userID,
			CategoryID: categoryID,
		}
		createdTag, err = s.TagRepository.Create(ctx, alt)
		if err != nil {
			return 0, fmt.Errorf("create category tag: %w", err)
		}
	}

	return createdTag.ID, nil
}

// createAutoTag creates a new auto-generated tag for a category.
func (s *Service) createAutoTag(ctx context.Context, categoryID uint64, userID uint64) (uint64, error) {
	newTag := tagdomain.Tag{
		Name:       fmt.Sprintf("auto-tag-%d", time.Now().UTC().UnixNano()),
		UserID:     userID,
		CategoryID: categoryID,
	}

	createdTag, err := s.TagRepository.Create(ctx, newTag)
	if err != nil {
		return 0, fmt.Errorf("create auto tag: %w", err)
	}

	return createdTag.ID, nil
}
