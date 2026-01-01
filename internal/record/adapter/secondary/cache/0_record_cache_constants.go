// Package cache constants contains constants used for record cache operations.
package cache

import "time"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer for record cache operations.
const TracerName = "aionapi.record.cache"

// Span Names.
const (
	SpanNameRecordSave       = "record.cache.save"
	SpanNameRecordGet        = "record.cache.get"
	SpanNameRecordDelete     = "record.cache.delete"
	SpanNameRecordListSave   = "record.cache.list_save"
	SpanNameRecordListGet    = "record.cache.list_get"
	SpanNameRecordListDelete = "record.cache.list_delete"
)

// Operations.
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

// RecordExpirationDefault defines the default expiration for individual records (30 minutes).
const RecordExpirationDefault = 30 * time.Minute

// RecordListExpirationDefault defines the default expiration for record lists (15 minutes).
// Shorter than individual records due to high mutation rate.
const RecordListExpirationDefault = 15 * time.Minute

// Key formats.
const (
	// RecordIDKeyFormat : record:user:{userID}:id:{recordID}.
	RecordIDKeyFormat = "record:user:%d:id:%d"

	// RecordDayKeyFormat : record:user:{userID}:day:{YYYY-MM-DD}.
	RecordDayKeyFormat = "record:user:%d:day:%s"

	// RecordByCategoryKeyFormat : record:category:{categoryID}:user:{userID}.
	RecordByCategoryKeyFormat = "record:category:%d:user:%d"

	// RecordByTagKeyFormat : record:tag:{tagID}:user:{userID}.
	RecordByTagKeyFormat = "record:tag:%d:user:%d"
)

// Success messages.
const (
	RecordRetrievedSuccessfully     = "record retrieved successfully from cache"
	RecordDeletedSuccessfully       = "record deleted successfully from cache"
	RecordSavedSuccessfully         = "record saved successfully to cache" // #nosec G101 - false positive, this is a log message
	RecordListRetrievedSuccessfully = "record list retrieved successfully from cache"
	RecordListSavedSuccessfully     = "record list saved successfully to cache" // #nosec G101 - false positive, this is a log message
	RecordListDeletedSuccessfully   = "record list deleted successfully from cache"
)

// Error messages.
const (
	ErrorToSaveRecordToCache         = "error to save record to cache"
	ErrorToGetRecordFromCache        = "error to get record from cache"
	ErrorToDeleteRecordFromCache     = "error to delete record from cache"
	ErrorToSerializeRecord           = "error to serialize record"
	ErrorToDeserializeRecord         = "error to deserialize record"
	ErrorToSaveRecordListToCache     = "error to save record list to cache"
	ErrorToGetRecordListFromCache    = "error to get record list from cache"
	ErrorToDeleteRecordListFromCache = "error to delete record list from cache"
	ErrorToSerializeRecordList       = "error to serialize record list"
	ErrorToDeserializeRecordList     = "error to deserialize record list"
)
