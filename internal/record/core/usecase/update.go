package usecase

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
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

	// Apply patch-like updates
	if cmd.Title != nil {
		existing.Title = *cmd.Title
	}
	if cmd.Description != nil {
		existing.Description = cmd.Description
	}
	if cmd.CategoryID != nil {
		existing.CategoryID = *cmd.CategoryID
	}
	if cmd.TagID != nil {
		existing.TagID = *cmd.TagID
	}
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
