// Package controller contains GraphQL-facing controllers for the Record context.
package controller

import "errors"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName is the name of the tracer for record GraphQL controllers.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.record.controller"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanCreate is the span name for creating a record.
	SpanCreate = "record.controller.create"

	// SpanGetByID is the span name for retrieving a record by ID.
	SpanGetByID = "record.controller.get_by_id"

	// SpanGetProjectedByID is the span name for retrieving a derived record projection by ID.
	SpanGetProjectedByID = "record.controller.get_projected_by_id"

	// SpanGetByName is the span name for retrieving a record by name.
	SpanGetByName = "record.controller.get_by_name"

	// SpanGetByCategory is the span name for retrieving records by category.
	SpanGetByCategory = "record.controller.get_by_category"

	// SpanGetByTag is the span name for retrieving a record by tag.
	SpanGetByTag = "record.controller.get_by_tag"

	// SpanListAll is the span name for listing all records for a user.
	SpanListAll = "record.controller.list_all"

	// SpanListByTag is the span name for listing records by tag.
	SpanListByTag = "record.controller.list_by_tag"

	// SpanListByCategory is the span name for listing records by category.
	SpanListByCategory = "record.controller.list_by_category"

	// SpanListLatest is the span name for listing latest records.
	SpanListLatest = "record.controller.list_latest"

	// SpanListProjectedLatest is the span name for listing latest derived record projections.
	SpanListProjectedLatest = "record.controller.list_projected_latest"

	// SpanListProjectedPage is the span name for cursor-based derived projection listing.
	SpanListProjectedPage = "record.controller.list_projected_page"

	// SpanListByDay is the span name for listing records by day.
	SpanListByDay = "record.controller.list_by_day"

	// SpanListAllUntil is the span name for listing records until a timestamp.
	SpanListAllUntil = "record.controller.list_all_until"

	// SpanListAllBetween is the span name for listing records between dates.
	SpanListAllBetween = "record.controller.list_all_between"

	// SpanRecordStats is the span name for aggregated record stats.
	SpanRecordStats = "record.controller.record_stats"

	// SpanSearch is the span name for searching records.
	SpanSearch = "record.controller.search"

	// SpanUpdate is the span name for updating a record.
	SpanUpdate = "record.controller.update"

	// SpanSoftDelete is the span name for soft-deleting a record.
	SpanSoftDelete = "record.controller.soft_delete"

	// SpanSoftDeleteAll is the span name for soft-deleting all records.
	SpanSoftDeleteAll = "record.controller.soft_delete_all"

	// SpanInsightFeed is the span name for insight feed queries.
	SpanInsightFeed = "record.controller.insight_feed"

	// SpanAnalyticsSeries is the span name for analytics series queries.
	SpanAnalyticsSeries = "record.controller.analytics_series"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// MsgCreated is the log message for when a record is created.
	MsgCreated = "record created"

	// StatusCreated is the human-readable status used in responses/logs when a record is created.
	StatusCreated = "created"

	// StatusFetched is the human-readable status used in responses/logs when records are retrieved.
	StatusFetched = "fetched"

	// StatusSearchCompleted is the human-readable status for search completion.
	StatusSearchCompleted = "search completed"

	// StatusStatsComputed is the human-readable status for stats computation.
	StatusStatsComputed = "record stats computed"

	// StatusDeleted is the human-readable status for soft deletion.
	StatusDeleted = "deleted"

	// StatusDeletedAll is the human-readable status for soft deletion of all records.
	StatusDeletedAll = "deleted_all"

	// StatusUpdated is the human-readable status for record update.
	StatusUpdated = "updated"
)

// =============================================================================
// BUSINESS LOGIC - Error Messages
// =============================================================================

// -----------------------------------------------------------------------------
// Log Messages
// -----------------------------------------------------------------------------

const (
	// MsgCreateError is the log message for when a create operation fails.
	MsgCreateError = "error creating record"

	// MsgSearchError is the log message for when a search operation fails.
	MsgSearchError = "error searching records"

	// MsgSearched is the log message for when records are searched successfully.
	MsgSearched = "records searched"

	// MsgStatsError is the log message for when a stats operation fails.
	MsgStatsError = "error computing record stats"

	// MsgListByTagError is the log message for when a list by tag operation fails.
	MsgListByTagError = "error listing records by tag"

	// MsgListByCategoryError is the log message for when a list by category operation fails.
	MsgListByCategoryError = "error listing records by category"

	// MsgListByDayError is the log message for when a list by day operation fails.
	MsgListByDayError = "error listing records by day"

	// MsgListUntilError is the log message for when a list until operation fails.
	MsgListUntilError = "error listing records until"

	// MsgListBetweenError is the log message for when a list between operation fails.
	MsgListBetweenError = "error listing records between dates"

	// MsgListLatestError is the log message for when a list latest operation fails.
	MsgListLatestError = "error listing latest records"

	// MsgGetProjectedByIDError is the log message for projected record retrieval failures.
	MsgGetProjectedByIDError = "error getting projected record by id"

	// MsgListProjectedLatestError is the log message for listing projected records failures.
	MsgListProjectedLatestError = "error listing latest projected records"

	// MsgListProjectedPageError is the log message for cursor-based projected records failures.
	MsgListProjectedPageError = "error listing projected records page"

	// MsgInvalidDateFormat is the log message for invalid date format.
	MsgInvalidDateFormat = "invalid date format"

	// MsgInvalidUntilTimestamp is the log message for invalid until timestamp.
	MsgInvalidUntilTimestamp = "invalid until timestamp"

	// MsgInvalidStartDate is the log message for invalid start date.
	MsgInvalidStartDate = "invalid start date"

	// MsgInvalidEndDate is the log message for invalid end date.
	MsgInvalidEndDate = "invalid end date"

	// MsgSoftDeleteError is the log message for soft delete operation failure.
	MsgSoftDeleteError = "error soft deleting record"

	// MsgSoftDeleteAllError is the log message for soft delete all operation failure.
	MsgSoftDeleteAllError = "error soft deleting all records"

	// MsgUpdateError is the log message for update operation failure.
	MsgUpdateError = "error updating record"

	// MsgInsightFeedError is the log message for insight feed failures.
	MsgInsightFeedError = "error computing insight feed"

	// MsgAnalyticsSeriesError is the log message for analytics series failures.
	MsgAnalyticsSeriesError = "error computing analytics series"

	// MsgInsightFeedFetched is the log message for successful insight queries.
	MsgInsightFeedFetched = "insight feed fetched"

	// MsgAnalyticsSeriesComputed is the log message for successful analytics series queries.
	MsgAnalyticsSeriesComputed = "analytics series computed"
)

// -----------------------------------------------------------------------------
// Attribute Keys
// Used in span.SetAttributes() and log metadata
// -----------------------------------------------------------------------------

const (
	// AttrResultsCount is the attribute key for result count.
	AttrResultsCount = "results_count"

	// AttrLimit is the attribute key for pagination limit.
	AttrLimit = "limit"

	// AttrCount is the attribute key for general count.
	AttrCount = "count"

	// AttrDate is the attribute key for date.
	AttrDate = "date"

	// AttrStartDate is the attribute key for start date.
	AttrStartDate = "start_date"

	// AttrEndDate is the attribute key for end date.
	AttrEndDate = "end_date"

	// AttrUntil is the attribute key for until timestamp.
	AttrUntil = "until"

	// AttrRecordsCount is the attribute key for records count.
	AttrRecordsCount = "records_count"

	// AttrTimezone is the attribute key for timezone.
	AttrTimezone = "timezone"

	// AttrSeriesKey is the attribute key for analytics series key.
	AttrSeriesKey = "series_key"

	// AttrWindow is the attribute key for insight window.
	AttrWindow = "window"

	// AttrTagIDsCount is the attribute key for scoped tag ids count.
	AttrTagIDsCount = "tag_ids_count"
)

// -----------------------------------------------------------------------------
// Date Parsing Defaults
// -----------------------------------------------------------------------------

const (
	// DefaultQueryTimezone is the default timezone for date parsing queries.
	DefaultQueryTimezone = "America/Sao_Paulo"
)

// -----------------------------------------------------------------------------
// Sentinel Errors
// Use errors.Is() for type-safe error comparison
// -----------------------------------------------------------------------------

var (
	// ErrUserIDNotFound is the error when the user ID is missing or invalid.
	ErrUserIDNotFound = errors.New("user id not found")

	// ErrRecordNotFound is the error when the record ID is missing or invalid.
	ErrRecordNotFound = errors.New("record id not found")

	// ErrInvalidRecordID is the error when the record ID cannot be parsed or is invalid.
	ErrInvalidRecordID = errors.New("invalid record id")

	// ErrInvalidCategoryID is the error when the category ID cannot be parsed or is invalid.
	ErrInvalidCategoryID = errors.New("invalid category id")

	// ErrInvalidTagID is the error when the tag ID cannot be parsed or is invalid.
	ErrInvalidTagID = errors.New("invalid tag id")

	// ErrTagNotFound is the error when a tag cannot be found.
	ErrTagNotFound = errors.New("tag not found")

	// ErrFailedToListRecords is the error when listing records fails.
	ErrFailedToListRecords = errors.New("failed to list records")

	// ErrInvalidRecordQueryDate is returned when a natural-language date cannot be parsed.
	ErrInvalidRecordQueryDate = errors.New(
		"invalid date format, expected YYYY-MM-DD, DD/MM/YYYY, RFC3339, or natural language (today, yesterday, ...)",
	)

	// ErrTagIDCannotBeZero is returned when a tag ID is zero.
	ErrTagIDCannotBeZero = errors.New("tag id cannot be zero")

	// ErrCategoryIDCannotBeZero is returned when a category ID is zero.
	ErrCategoryIDCannotBeZero = errors.New("category id cannot be zero")

	// ErrInvalidUntilTimestamp is returned when the until timestamp is invalid.
	ErrInvalidUntilTimestamp = errors.New("invalid until timestamp, expected RFC3339")

	// ErrInvalidStartDate is returned when the start date is invalid.
	ErrInvalidStartDate = errors.New("invalid start date, expected RFC3339")

	// ErrInvalidEndDate is returned when the end date is invalid.
	ErrInvalidEndDate = errors.New("invalid end date, expected RFC3339")
)
