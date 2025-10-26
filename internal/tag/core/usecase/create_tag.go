package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"github.com/lechitz/AionApi/internal/tag/core/ports/input"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Create creates a new Tag after validating inputs and ensuring name uniqueness.
// It receives a CreateTagCommand (input port), builds the domain entity and delegates to the repository.
func (s *Service) Create(ctx context.Context, cmd input.CreateTagCommand) (domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanCreateTag)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanCreateTag),
		attribute.String(commonkeys.TagName, cmd.Name),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(cmd.CategoryID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(cmd.UserID, 10)),
	)

	span.AddEvent(EventValidateInput)
	if err := s.validateCreateCommand(cmd); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrToValidateTag)
		s.Logger.ErrorwCtx(ctx, ErrToValidateTag, commonkeys.Error, err.Error())
		return domain.Tag{}, err
	}

	// Build domain entity from the command (no GraphQL types here).
	newTag := domain.Tag{
		UserID:      cmd.UserID,
		CategoryID:  cmd.CategoryID,
		Name:        cmd.Name,
		Description: ptrOrEmpty(cmd.Description),
	}

	span.AddEvent(EventCheckUniqueness)
	existingCategory, err := s.TagRepository.GetByName(ctx, newTag.Name, newTag.UserID)
	if err == nil && existingCategory.Name != "" {
		span.SetStatus(codes.Error, TagAlreadyExists)
		s.Logger.ErrorwCtx(ctx, TagAlreadyExists, commonkeys.CategoryName, newTag.Name)
		return domain.Tag{}, errors.New(TagAlreadyExists)
	}

	span.AddEvent(EventRepositoryCreate)
	createdTag, err := s.TagRepository.Create(ctx, newTag)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToCreateTag)
		s.Logger.ErrorwCtx(ctx, FailedToCreateTag, commonkeys.Tag, newTag, commonkeys.Error, err)
		return domain.Tag{}, fmt.Errorf("%s: %w", FailedToCreateTag, err)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusCreated)
	s.Logger.InfowCtx(ctx, fmt.Sprintf(SuccessfullyCreatedTag, createdTag.Name))

	return createdTag, nil
}

// validateCreateCommand checks required fields and length constraints for CreateCategoryCommand.
func (s *Service) validateCreateCommand(cmd input.CreateTagCommand) error {
	if cmd.UserID == 0 {
		return errors.New(UserIDIsRequired)
	}
	if cmd.Name == "" {
		return errors.New(TagNameIsRequired)
	}
	if cmd.Description != nil && len(*cmd.Description) > 200 {
		return errors.New(TagDescriptionIsTooLong)
	}
	return nil
}

// ptrOrEmpty converts a *string into a safe string, returning "" when nil.
func ptrOrEmpty(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
