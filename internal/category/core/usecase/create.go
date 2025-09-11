package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Create creates a new handler in the database after validating inputs and ensuring uniqueness by name.
// Returns the created handler or an error.
func (s *Service) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanCreateCategory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanCreateCategory),
		attribute.String(commonkeys.CategoryName, category.Name),
	)

	span.AddEvent(EventValidateInput)
	if err := s.validateCreateCategoryRequired(category); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrToValidateCategory)
		s.Logger.ErrorwCtx(ctx, ErrToValidateCategory, commonkeys.Error, err.Error())
		return domain.Category{}, err
	}

	span.AddEvent(EventCheckUniqueness)
	existingCategory, err := s.CategoryRepository.GetByName(ctx, category.Name, category.UserID)
	if err == nil && existingCategory.Name != "" {
		span.SetStatus(codes.Error, CategoryAlreadyExists)
		s.Logger.ErrorwCtx(ctx, CategoryAlreadyExists, commonkeys.CategoryName, category.Name)
		return domain.Category{}, errors.New(CategoryAlreadyExists)
	}

	span.AddEvent(EventRepositoryCreate)
	createdCategory, err := s.CategoryRepository.Create(ctx, category)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToCreateCategory)
		s.Logger.ErrorwCtx(ctx, FailedToCreateCategory, commonkeys.Category, category, commonkeys.Error, err)
		return domain.Category{}, fmt.Errorf("%s: %w", FailedToCreateCategory, err)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, "created")
	s.Logger.InfowCtx(ctx, fmt.Sprintf(SuccessfullyCreatedCategory, category.Name))

	return createdCategory, nil
}

// validateCreateCategoryRequired validates required fields for creating a handler and enforces constraints like name presence and field length limits.
func (s *Service) validateCreateCategoryRequired(category domain.Category) error {
	if category.Name == "" {
		return errors.New(CategoryNameIsRequired)
	}

	if category.Description != "" && len(category.Description) > 200 {
		return errors.New(CategoryDescriptionIsTooLong)
	}

	if category.Color != "" && len(category.Color) > 7 {
		return errors.New(CategoryColorIsTooLong)
	}

	if category.Icon != "" && len(category.Icon) > 50 {
		return errors.New(CategoryIconIsTooLong)
	}

	return nil
}
