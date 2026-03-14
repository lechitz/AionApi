// Package cache constants contains constants used for category cache operations.
package cache

import "time"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer for category cache operations.
// Format: aionapi.<domain>.<layer>.
const TracerName = "aionapi.category.cache"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanNameCategorySave is the span name for saving a category.
	SpanNameCategorySave = "category.cache.save"

	// SpanNameCategoryGet is the span name for retrieving a category.
	SpanNameCategoryGet = "category.cache.get"

	// SpanNameCategoryDelete is the span name for deleting a category.
	SpanNameCategoryDelete = "category.cache.delete"

	// SpanNameCategoryListSave is the span name for saving category list.
	SpanNameCategoryListSave = "category.cache.list_save"

	// SpanNameCategoryListGet is the span name for retrieving category list.
	SpanNameCategoryListGet = "category.cache.list_get"

	// SpanNameCategoryListDelete is the span name for deleting category list.
	SpanNameCategoryListDelete = "category.cache.list_delete"
)

// -----------------------------------------------------------------------------
// Span Attributes.
// -----------------------------------------------------------------------------

const (
	// OperationSave is the name of the save operation.
	OperationSave = "save"

	// OperationGet is the name of the get operation.
	OperationGet = "get"

	// OperationDelete is the name of the delete operation.
	OperationDelete = "delete"
)

// Attribute keys for tracing and logging.
const (
	// AttributeCacheKey is the attribute key for cache keys.
	AttributeCacheKey = "cache_key"

	// AttributeTTL is the attribute key for time-to-live values.
	AttributeTTL = "ttl"
)

// =============================================================================
// BUSINESS LOGIC - Configuration and Messages.
// =============================================================================

// CategoryExpirationDefault defines the default expiration for categories (1 hour).
const CategoryExpirationDefault = 1 * time.Hour

// CategoryListExpirationDefault defines the default expiration for category lists (1 hour).
const CategoryListExpirationDefault = 1 * time.Hour

// CategoryIDKeyFormat defines the format for category keys by ID.
// Format: category:user:{userID}:id:{categoryID}.
const CategoryIDKeyFormat = "category:user:%d:id:%d"

// CategoryNameKeyFormat defines the format for category keys by name.
// Format: category:user:{userID}:name:{categoryName}.
const CategoryNameKeyFormat = "category:user:%d:name:%s"

// CategoryListKeyFormat defines the format for category list keys.
// Format: category:user:{userID}:list.
const CategoryListKeyFormat = "category:user:%d:list"

// Success messages.
const (
	// CategoryRetrievedSuccessfully is the message used when a category is retrieved successfully.
	CategoryRetrievedSuccessfully = "category retrieved successfully from cache"

	// CategoryDeletedSuccessfully is the message used when a category is deleted successfully.
	CategoryDeletedSuccessfully = "category deleted successfully from cache"

	// CategorySavedSuccessfully is the message used when a category is saved successfully.
	CategorySavedSuccessfully = "category saved successfully to cache" // #nosec G101 - false positive, this is a log message

	// CategoryListRetrievedSuccessfully is the message used when a category list is retrieved successfully.
	CategoryListRetrievedSuccessfully = "category list retrieved successfully from cache"

	// CategoryListSavedSuccessfully is the message used when a category list is saved successfully.
	CategoryListSavedSuccessfully = "category list saved successfully to cache" // #nosec G101 - false positive, this is a log message

	// CategoryListDeletedSuccessfully is the message used when a category list is deleted successfully.
	CategoryListDeletedSuccessfully = "category list deleted successfully from cache"
)

// Error messages.
const (
	// ErrorToSaveCategoryToCache indicates a failure to save a category in cache.
	ErrorToSaveCategoryToCache = "error to save category to cache"

	// ErrorToGetCategoryFromCache indicates a failure to retrieve a category from cache.
	ErrorToGetCategoryFromCache = "error to get category from cache"

	// ErrorToDeleteCategoryFromCache indicates a failure to delete a category from cache.
	ErrorToDeleteCategoryFromCache = "error to delete category from cache"

	// ErrorToSerializeCategory indicates a failure to serialize a category.
	ErrorToSerializeCategory = "error to serialize category"

	// ErrorToDeserializeCategory indicates a failure to deserialize a category.
	ErrorToDeserializeCategory = "error to deserialize category"

	// ErrorToSaveCategoryListToCache indicates a failure to save a category list in cache.
	ErrorToSaveCategoryListToCache = "error to save category list to cache"

	// ErrorToGetCategoryListFromCache indicates a failure to retrieve a category list from cache.
	ErrorToGetCategoryListFromCache = "error to get category list from cache"

	// ErrorToDeleteCategoryListFromCache indicates a failure to delete a category list from cache.
	ErrorToDeleteCategoryListFromCache = "error to delete category list from cache"

	// ErrorToSerializeCategoryList indicates a failure to serialize a category list.
	ErrorToSerializeCategoryList = "error to serialize category list"

	// ErrorToDeserializeCategoryList indicates a failure to deserialize a category list.
	ErrorToDeserializeCategoryList = "error to deserialize category list"
)
