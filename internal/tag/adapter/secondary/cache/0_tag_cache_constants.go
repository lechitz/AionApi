// Package cache constants contains constants used for tag cache operations.
package cache

import "time"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer for tag cache operations.
const TracerName = "aionapi.tag.cache"

// Span Names
const (
	SpanNameTagSave       = "tag.cache.save"
	SpanNameTagGet        = "tag.cache.get"
	SpanNameTagDelete     = "tag.cache.delete"
	SpanNameTagListSave   = "tag.cache.list_save"
	SpanNameTagListGet    = "tag.cache.list_get"
	SpanNameTagListDelete = "tag.cache.list_delete"
)

// Operations
const (
	OperationSave   = "save"
	OperationGet    = "get"
	OperationDelete = "delete"
)

// Attribute keys for tracing and logging.
const (
	AttributeCacheKey = "cache_key"
	AttributeTTL      = "ttl"
)

// =============================================================================
// BUSINESS LOGIC - Configuration and Messages
// =============================================================================

// TagExpirationDefault defines the default expiration for tags (1 hour).
const TagExpirationDefault = 1 * time.Hour

// TagListExpirationDefault defines the default expiration for tag lists (1 hour).
const TagListExpirationDefault = 1 * time.Hour

// Key formats
const (
	// TagIDKeyFormat: tag:user:{userID}:id:{tagID}
	TagIDKeyFormat = "tag:user:%d:id:%d"

	// TagNameKeyFormat: tag:user:{userID}:name:{tagName}
	TagNameKeyFormat = "tag:user:%d:name:%s"

	// TagListKeyFormat: tag:user:{userID}:list
	TagListKeyFormat = "tag:user:%d:list"

	// TagByCategoryKeyFormat: tag:category:{categoryID}:user:{userID}
	TagByCategoryKeyFormat = "tag:category:%d:user:%d"
)

// Success messages.
const (
	TagRetrievedSuccessfully     = "tag retrieved successfully from cache"
	TagDeletedSuccessfully       = "tag deleted successfully from cache"
	TagSavedSuccessfully         = "tag saved successfully to cache"
	TagListRetrievedSuccessfully = "tag list retrieved successfully from cache"
	TagListSavedSuccessfully     = "tag list saved successfully to cache"
	TagListDeletedSuccessfully   = "tag list deleted successfully from cache"
)

// Error messages.
const (
	ErrorToSaveTagToCache         = "error to save tag to cache"
	ErrorToGetTagFromCache        = "error to get tag from cache"
	ErrorToDeleteTagFromCache     = "error to delete tag from cache"
	ErrorToSerializeTag           = "error to serialize tag"
	ErrorToDeserializeTag         = "error to deserialize tag"
	ErrorToSaveTagListToCache     = "error to save tag list to cache"
	ErrorToGetTagListFromCache    = "error to get tag list from cache"
	ErrorToDeleteTagListFromCache = "error to delete tag list from cache"
	ErrorToSerializeTagList       = "error to serialize tag list"
	ErrorToDeserializeTagList     = "error to deserialize tag list"
)
