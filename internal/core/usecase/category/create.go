package category

import (
	"context"
	"errors"
	"fmt"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
)

// Create creates a new category in the database after validating inputs and ensuring uniqueness by name.
// Returns the created category or an error.
func (s *Service) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	tr := otel.Tracer(constants.TracerName)
	ctx, span := tr.Start(ctx, constants.SpanCreateCategory)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanCreateCategory),
		attribute.String(commonkeys.CategoryName, category.Name),
	)

	span.AddEvent(constants.EventValidateInput)
	if err := s.validateCreateCategoryRequired(category); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrToValidateCategory)
		s.Logger.ErrorwCtx(ctx, constants.ErrToValidateCategory, commonkeys.Error, err.Error())
		return domain.Category{}, err
	}

	span.AddEvent(constants.EventCheckUniqueness)
	existingCategory, err := s.CategoryRepository.GetByName(ctx, category.Name, category.UserID)
	if err == nil && existingCategory.Name != "" {
		span.SetStatus(codes.Error, constants.CategoryAlreadyExists)
		s.Logger.ErrorwCtx(ctx, constants.CategoryAlreadyExists, commonkeys.CategoryName, category.Name)
		return domain.Category{}, errors.New(constants.CategoryAlreadyExists)
	}

	span.AddEvent(constants.EventRepositoryCreate)
	createdCategory, err := s.CategoryRepository.Create(ctx, category)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.FailedToCreateCategory)
		s.Logger.ErrorwCtx(ctx, constants.FailedToCreateCategory, commonkeys.Category, category, commonkeys.Error, err)
		return domain.Category{}, fmt.Errorf("%s: %w", constants.FailedToCreateCategory, err)
	}

	span.AddEvent(constants.EventSuccess)
	span.SetStatus(codes.Ok, "created")
	s.Logger.InfowCtx(ctx, fmt.Sprintf(constants.SuccessfullyCreatedCategory, category.Name))

	return createdCategory, nil
}

// validateCreateCategoryRequired validates required fields for creating a category and enforces constraints like name presence and field length limits.
func (s *Service) validateCreateCategoryRequired(category domain.Category) error {
	if category.Name == "" {
		return errors.New(constants.CategoryNameIsRequired)
	}

	if category.Description != "" && len(category.Description) > 200 {
		return errors.New(constants.CategoryDescriptionIsTooLong)
	}

	if category.Color != "" && len(category.Color) > 7 {
		return errors.New(constants.CategoryColorIsTooLong)
	}

	if category.Icon != "" && len(category.Icon) > 50 {
		return errors.New(constants.CategoryIconIsTooLong)
	}

	return nil
}
