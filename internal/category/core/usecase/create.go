package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/category/core/ports/input"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Create creates a new Category after validating inputs and ensuring name uniqueness.
// It receives a CreateCategoryCommand (input port), builds the domain entity and delegates to the repository.
func (s *Service) Create(ctx context.Context, cmd input.CreateCategoryCommand) (domain.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanCreateCategory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanCreateCategory),
		attribute.String(commonkeys.CategoryName, cmd.Name),
	)

	span.AddEvent(EventValidateInput)
	if err := s.validateCreateCommand(cmd); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrToValidateCategory)
		s.Logger.ErrorwCtx(ctx, ErrToValidateCategory, commonkeys.Error, err.Error())
		return domain.Category{}, err
	}

	// Build domain entity from the command (no GraphQL types here).
	newCategory := domain.Category{
		UserID:      cmd.UserID,
		Name:        cmd.Name,
		Description: ptrOrEmpty(cmd.Description),
		Color:       ptrOrEmpty(cmd.ColorHex),
		Icon:        ptrOrEmpty(cmd.Icon),
	}

	span.AddEvent(EventCheckUniqueness)
	existingCategory, err := s.CategoryRepository.GetByName(ctx, newCategory.Name, newCategory.UserID)
	if err == nil && existingCategory.Name != "" {
		span.SetStatus(codes.Error, CategoryAlreadyExists)
		s.Logger.ErrorwCtx(ctx, CategoryAlreadyExists, commonkeys.CategoryName, newCategory.Name)
		return domain.Category{}, errors.New(CategoryAlreadyExists)
	}

	span.AddEvent(EventRepositoryCreate)
	createdCategory, err := s.CategoryRepository.Create(ctx, newCategory)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToCreateCategory)
		s.Logger.ErrorwCtx(ctx, FailedToCreateCategory, commonkeys.Category, newCategory, commonkeys.Error, err)
		return domain.Category{}, fmt.Errorf("%s: %w", FailedToCreateCategory, err)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, "created")
	s.Logger.InfowCtx(ctx, fmt.Sprintf(SuccessfullyCreatedCategory, createdCategory.Name))

	return createdCategory, nil
}

// validateCreateCommand checks required fields and length constraints for CreateCategoryCommand.
func (s *Service) validateCreateCommand(cmd input.CreateCategoryCommand) error {
	if cmd.Name == "" {
		return errors.New(CategoryNameIsRequired)
	}
	if cmd.Description != nil && len(*cmd.Description) > 200 {
		return errors.New(CategoryDescriptionIsTooLong)
	}
	if cmd.ColorHex != nil && len(*cmd.ColorHex) > 7 {
		return errors.New(CategoryColorIsTooLong)
	}
	if cmd.Icon != nil && len(*cmd.Icon) > 50 {
		return errors.New(CategoryIconIsTooLong)
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
