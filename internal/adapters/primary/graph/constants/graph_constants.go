// Package constants contains constants used throughout the application.
package constants

// UserID is the key for the user ID in context or responses.
const UserID = "user_id"

// CategoryID is the key for the category ID in context or responses.
const CategoryID = "category_id"

// CategoryName is the key for the category Name in context or responses.
const CategoryName = "category_name"

// CategoryColor is the key for the category Color in context or responses.
const CategoryColor = "category_color"

// CategoriesCount is the key for the number of categories in context or responses.
const CategoriesCount = "categories_count"

// Error is the key for error messages in context.
const Error = "error"

// ErrUserIDNotFound is the error message when the user ID is not found in the context.
const ErrUserIDNotFound = "user id not found in context"

// ErrCategoryNotFound is the error message when the category is not found.
const ErrCategoryNotFound = "category not found"

// InvalidCategoryID is the error message when the category ID is invalid.
const InvalidCategoryID = "invalid category id"

// ErrCategoryByNameNotFound is the error message when the name does not find the category.
const ErrCategoryByNameNotFound = "category not found by name"

// ErrAllCategoriesNotFound is the error message when no categories are found.
const ErrAllCategoriesNotFound = "no categories found"

// SuccessCategoryCreated is the success message when a category is created.
const SuccessCategoryCreated = "category created successfully"

// SuccessCategoryFetch is the success message when a category is fetch.
const SuccessCategoryFetch = "category fetched successfully"

// SuccessAllCategoriesFetch is the success message when all categories are fetch.
const SuccessAllCategoriesFetch = "all categories fetched successfully"

// SuccessCategoryUpdated is the success message when the category is updated.
const SuccessCategoryUpdated = "category updated successfully"

// SuccessCategorySoftDeleted is the success message when the category is soft deleted.
const SuccessCategorySoftDeleted = "category soft deleted successfully"

// SpanTracerCategory is the name of the tracer category for GraphQL.
const SpanTracerCategory = "GraphQL/Category"

// SpanStartCreateCategory is the name of the span for the CreateCategory mutation.
const SpanStartCreateCategory = "CreateCategory"

// SpanEventCreateCategory is the name of the event for the CreateCategory mutation.
const SpanEventCreateCategory = "start CreateCategory mutation"

// SpanEventGetCategoryByID is the start span for the GetCategoryByID query.
const SpanEventGetCategoryByID = "start GetCategoryByID query"

// SpanEventGetCategoryByName is the start span for the GetCategoryByName query.
const SpanEventGetCategoryByName = "start GetCategoryByName query"

// SpanEventUpdateCategory is the start span for the UpdateCategory mutation.
const SpanEventUpdateCategory = "start UpdateCategory mutation"

// SpanEventGetAllCategories is the start span for the GetAllCategories query.
const SpanEventGetAllCategories = "start GetAllCategories query"

// SpanEventSoftDeleteCategory is the start span for the SoftDeleteCategory mutation.
const SpanEventSoftDeleteCategory = "start SoftDeleteCategory mutation"

// SpanStartGetCategoryByID is the span name for the GetCategoryByID query.
const SpanStartGetCategoryByID = "GetCategoryByID"

// SpanStartAllGetCategories is the span name for the GetAllCategories query.
const SpanStartAllGetCategories = "GetAllCategories"

// SpanStartGetCategoryByName is the span name for the GetCategoryByName query.
const SpanStartGetCategoryByName = "GetCategoryByName"

// SpanStartUpdateCategory is the span name for the UpdateCategory mutation.
const SpanStartUpdateCategory = "UpdateCategory"

// SpanStartSoftDeleteCategory is the span name for the SoftDeleteCategory mutation.
const SpanStartSoftDeleteCategory = "SoftDeleteCategory"
