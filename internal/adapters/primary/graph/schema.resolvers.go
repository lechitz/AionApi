// Package graph implements the GraphQL resolvers for the API.
//
//nolint:govet,revive,perfsprint,nolintlint
package graph

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/lechitz/AionApi/internal/adapters/primary/graph/model"
	"github.com/lechitz/AionApi/internal/core/domain"
)

// traceAttributesFromCategory creates a slice of attributes from the given category.
func traceAttributesFromCategory(category model.DtoCreateCategory) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String("category_name", category.Name),
	}
	if category.Description != nil {
		attrs = append(attrs, attribute.String("category_description", *category.Description))
	}
	if category.ColorHex != nil {
		attrs = append(attrs, attribute.String("category_color", *category.ColorHex))
	}
	if category.Icon != nil {
		attrs = append(attrs, attribute.String("category_icon", *category.Icon))
	}
	return attrs
}

// CreateCategory is the resolver for the createCategory field.
func (r *mutationResolver) CreateCategory(ctx context.Context, category model.DtoCreateCategory) (*model.Category, error) {
	tracer := otel.Tracer("AionApi/GraphQL/Category")
	ctx, span := tracer.Start(ctx, "CreateCategoryResolver")
	defer span.End()

	span.AddEvent(
		"start createCategory mutation",
		trace.WithAttributes(traceAttributesFromCategory(category)...),
	)

	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		err := errors.New("userID not found in context")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	createCategory := domain.Category{
		UserID:      userID,
		Name:        category.Name,
		Description: *category.Description,
		Color:       *category.ColorHex,
		Icon:        *category.Icon,
	}

	categoryDB, err := r.CategoryService.CreateCategory(ctx, createCategory)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(
		attribute.String("category_name", categoryDB.Name),
		attribute.String("user_id", strconv.FormatUint(categoryDB.UserID, 10)),
		attribute.String("category_id", strconv.FormatUint(categoryDB.ID, 10)),
	)
	span.SetStatus(codes.Ok, "category created successfully")

	return &model.Category{
		CategoryID:  strconv.FormatUint(categoryDB.ID, 10),
		UserID:      strconv.FormatUint(userID, 10),
		Name:        categoryDB.Name,
		Description: &categoryDB.Description,
		ColorHex:    &categoryDB.Color,
		Icon:        &categoryDB.Icon,
	}, nil
}

// CreateTag is the resolver for the createTag field (stub for interface compliance).
func (r *mutationResolver) CreateTag(ctx context.Context, input model.NewTag) (*model.Tags, error) {
	return nil, errors.New("not implemented")
}

// UpdateCategory is the resolver for the updateCategory field.
func (r *mutationResolver) UpdateCategory(
	ctx context.Context,
	category model.DtoUpdateCategory,
) (*model.Category, error) {
	tracer := otel.Tracer("AionApi/GraphQL/Category")
	ctx, span := tracer.Start(ctx, "UpdateCategoryResolver")
	defer span.End()

	span.AddEvent(
		"start updateCategory mutation",
		trace.WithAttributes(attribute.String("input_category_id", category.CategoryID)),
	)

	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		err := errors.New("userID not found in context")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	categoryIDUint, err := strconv.ParseUint(category.CategoryID, 10, 64)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.New("invalid category ID format")
	}

	span.SetAttributes(
		attribute.String("user_id", strconv.FormatUint(userID, 10)),
		attribute.String("category_id", category.CategoryID),
		attribute.String("operation", "update"),
	)

	updateCategory := domain.Category{
		ID:     categoryIDUint,
		UserID: userID,
	}

	if category.Name != nil {
		updateCategory.Name = *category.Name
	}
	if category.Description != nil {
		updateCategory.Description = *category.Description
	}
	if category.ColorHex != nil {
		updateCategory.Color = *category.ColorHex
	}
	if category.Icon != nil {
		updateCategory.Icon = *category.Icon
	}

	categoryDB, err := r.CategoryService.UpdateCategory(ctx, updateCategory)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.String("updated_name", categoryDB.Name))
	span.SetStatus(codes.Ok, "category updated successfully")

	return &model.Category{
		CategoryID:  strconv.FormatUint(categoryDB.ID, 10),
		UserID:      fmt.Sprintf("%d", categoryDB.UserID),
		Name:        categoryDB.Name,
		Description: &categoryDB.Description,
		ColorHex:    &categoryDB.Color,
		Icon:        &categoryDB.Icon,
	}, nil
}

// SoftDeleteCategory is the resolver for the softDeleteCategory field.
func (r *mutationResolver) SoftDeleteCategory(
	ctx context.Context,
	category model.DtoDeleteCategory,
) (bool, error) {
	tracer := otel.Tracer("AionApi/GraphQL/Category")
	ctx, span := tracer.Start(ctx, "SoftDeleteCategoryResolver")
	defer span.End()

	span.AddEvent(
		"start softDeleteCategory mutation",
		trace.WithAttributes(attribute.String("input_category_id", category.CategoryID)),
	)

	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		err := errors.New("userID not found in context")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return false, err
	}

	categoryIDUint, err := strconv.ParseUint(category.CategoryID, 10, 64)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return false, fmt.Errorf("invalid category ID format: %w", err)
	}

	span.SetAttributes(
		attribute.String("user_id", strconv.FormatUint(userID, 10)),
		attribute.String("category_id", category.CategoryID),
		attribute.String("operation", "soft_delete"),
	)

	categoryDomain := domain.Category{
		ID:     categoryIDUint,
		UserID: userID,
	}

	if err := r.CategoryService.SoftDeleteCategory(ctx, categoryDomain); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return false, err
	}

	span.SetStatus(codes.Ok, "soft delete successful")
	return true, nil
}

// --- Query Resolvers ---

// AllCategories is the resolver for the allCategories field.
func (r *queryResolver) AllCategories(ctx context.Context) ([]*model.Category, error) {
	tracer := otel.Tracer("AionApi/GraphQL/Category")
	ctx, span := tracer.Start(ctx, "AllCategoriesResolver")
	defer span.End()

	span.AddEvent(
		"start allCategories query",
	)

	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		err := errors.New("userID not found in context")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.Logger.Errorw("User ID not found in context", "error", err.Error())
		return nil, err
	}
	span.SetAttributes(attribute.String("user_id", strconv.FormatUint(userID, 10)))

	categoryDB, err := r.CategoryService.GetAllCategories(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to fetch categories")
		return nil, errors.New("failed to fetch categories")
	}

	categories := make([]*model.Category, len(categoryDB))
	for i, category := range categoryDB {
		categories[i] = &model.Category{
			CategoryID:  strconv.FormatUint(category.ID, 10),
			UserID:      strconv.FormatUint(userID, 10),
			Name:        category.Name,
			Description: &category.Description,
			ColorHex:    &category.Color,
			Icon:        &category.Icon,
		}
	}
	span.SetAttributes(attribute.Int("categories_count", len(categories)))
	span.SetStatus(codes.Ok, "categories fetched successfully")

	return categories, nil
}

// GetCategoryByID is the resolver for the getCategoryByID field.
func (r *queryResolver) GetCategoryByID(
	ctx context.Context,
	categoryRequest model.DtoGetCategoryByID,
) (*model.Category, error) {
	tracer := otel.Tracer("AionApi/GraphQL/Category")
	ctx, span := tracer.Start(ctx, "GetCategoryByIDResolver")
	defer span.End()

	span.AddEvent(
		"start getCategoryByID query",
		trace.WithAttributes(attribute.String("input_category_id", categoryRequest.CategoryID)),
	)

	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		err := errors.New("userID not found in context")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.Logger.Errorw("User ID not found in context", "error", err.Error())
		return nil, err
	}
	span.SetAttributes(attribute.String("user_id", strconv.FormatUint(userID, 10)))

	categoryIDUint, err := strconv.ParseUint(categoryRequest.CategoryID, 10, 64)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid category ID format")
		return nil, errors.New("invalid category ID format")
	}
	span.SetAttributes(attribute.String("category_id", categoryRequest.CategoryID))

	category := domain.Category{
		ID:     categoryIDUint,
		UserID: userID,
	}

	categoryDB, err := r.CategoryService.GetCategoryByID(ctx, category)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "category not found")
		return nil, err
	}

	span.SetStatus(codes.Ok, "category fetched successfully")
	span.SetAttributes(
		attribute.String("category_name", categoryDB.Name),
		attribute.String("category_color", categoryDB.Color),
	)

	return &model.Category{
		CategoryID:  strconv.FormatUint(categoryDB.ID, 10),
		UserID:      strconv.FormatUint(categoryDB.UserID, 10),
		Name:        categoryDB.Name,
		Description: &categoryDB.Description,
		ColorHex:    &categoryDB.Color,
		Icon:        &categoryDB.Icon,
	}, nil
}

// GetCategoryByName is the resolver for the getCategoryByName field.
func (r *queryResolver) GetCategoryByName(ctx context.Context, categoryRequest model.DtoGetCategoryByName) (*model.Category, error) {
	tracer := otel.Tracer("AionApi/GraphQL/Category")
	ctx, span := tracer.Start(ctx, "GetCategoryByNameResolver")
	defer span.End()

	span.AddEvent(
		"start getCategoryByName query",
		trace.WithAttributes(attribute.String("input_category_name", categoryRequest.Name)),
	)

	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		err := errors.New("userID not found in context")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		r.Logger.Errorw("User ID not found in context", "error", err.Error())
		return nil, err
	}
	span.SetAttributes(
		attribute.String("user_id", strconv.FormatUint(userID, 10)),
		attribute.String("category_name", categoryRequest.Name),
	)

	category := domain.Category{
		UserID: userID,
		Name:   categoryRequest.Name,
	}

	categoryDB, err := r.CategoryService.GetCategoryByName(ctx, category)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error getting category by name")
		return nil, err
	}
	if categoryDB.Name == "" {
		span.SetStatus(codes.Ok, "category not found (normal)")
		return nil, errors.New("category not found")
	}

	span.SetStatus(codes.Ok, "category fetched successfully")
	span.SetAttributes(
		attribute.String("category_id", fmt.Sprintf("%d", categoryDB.ID)),
		attribute.String("category_color", categoryDB.Color),
	)

	return &model.Category{
		CategoryID:  fmt.Sprintf("%d", categoryDB.ID),
		UserID:      fmt.Sprintf("%d", categoryDB.UserID),
		Name:        categoryDB.Name,
		Description: &categoryDB.Description,
		ColorHex:    &categoryDB.Color,
		Icon:        &categoryDB.Icon,
	}, nil
}

// GetAllTags is the resolver for the getAllTags field (not implemented).
func (r *queryResolver) GetAllTags(ctx context.Context) ([]*model.Tags, error) {
	panic(fmt.Errorf("not implemented: GetAllTags - GetAllTags"))
}

// GetTagByID is the resolver for the getTagByID field (not implemented).
func (r *queryResolver) GetTagByID(ctx context.Context, tagID string) (*model.Tags, error) {
	panic(fmt.Errorf("not implemented: GetTagByID - GetTagByID"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
