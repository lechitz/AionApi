package usecase

import "errors"

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerName identifies the OpenTelemetry tracer for the Record usecase package.
// Format: aionapi.<domain>.<layer> .
const TracerName = "aionapi.record.usecase"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanCreate is the span name for creating a record.
	SpanCreate = "record.create"

	// SpanGetByID is the span name for getting a record by ID.
	SpanGetByID = "record.get_by_id"

	// SpanGetByTag is the span name for getting records by tag.
	SpanGetByTag = "record.get_by_tag"

	// SpanGetByCategory is the span name for getting records by category.
	SpanGetByCategory = "record.get_by_category"

	// SpanListAll is the span name for listing all records.
	SpanListAll = "record.list_all"

	// SpanListByTag is the span name for listing records by tag.
	SpanListByTag = "record.list_by_tag"

	// SpanListByCategory is the span name for listing records by category.
	SpanListByCategory = "record.list_by_category"

	// SpanListLatest is the span name for listing latest records.
	SpanListLatest = "record.list_latest"

	// SpanListByDay is the span name for listing records by day.
	SpanListByDay = "record.list_by_day"

	// SpanListAllUntil is the span name for listing records until a timestamp.
	SpanListAllUntil = "record.list_all_until"

	// SpanListAllBetween is the span name for listing records between dates.
	SpanListAllBetween = "record.list_all_between"

	// SpanUpdate is the span name for updating a record.
	SpanUpdate = "record.update"

	// SpanSoftDelete is the span name for soft-deleting a record.
	SpanSoftDelete = "record.soft_delete"
)

// -----------------------------------------------------------------------------
// Event Names
// Format: <domain>.<action>.<detail>
// -----------------------------------------------------------------------------

const (
	// EventValidateInput marks the input-validation step.
	EventValidateInput = "record.input.validate"

	// EventRepositoryCreate marks the repository create call.
	EventRepositoryCreate = "record.repository.create"

	// EventRepositoryGet marks the repository single-get call.
	EventRepositoryGet = "record.repository.get"

	// EventRepositoryList marks the repository list call.
	EventRepositoryList = "record.repository.list"

	// EventRepositoryUpdate marks the repository update call.
	EventRepositoryUpdate = "record.repository.update"

	// EventRepositoryDelete marks the repository delete call.
	EventRepositoryDelete = "record.repository.delete"

	// EventCheckCache marks checking cache for records.
	EventCheckCache = "record.cache.check"

	// EventCacheHit marks a cache hit.
	EventCacheHit = "record.cache.hit"

	// EventCacheMiss marks a cache miss.
	EventCacheMiss = "record.cache.miss"

	// EventSaveToCache marks saving records to cache.
	EventSaveToCache = "record.cache.save"

	// EventSuccess marks a successful outcome.
	EventSuccess = "record.success"
)

// -----------------------------------------------------------------------------
// Status Descriptions
// -----------------------------------------------------------------------------

const (
	// StatusCreated indicates a record was created.
	StatusCreated = "created"

	// StatusRetrieved indicates a record was retrieved.
	StatusRetrieved = "retrieved"

	// StatusUpdated indicates a record was updated.
	StatusUpdated = "updated"

	// StatusDeleted indicates a record was deleted.
	StatusDeleted = "deleted"

	// StatusListedAll indicates all records were listed.
	StatusListedAll = "listed_all"
)

// =============================================================================
// BUSINESS LOGIC - Error and Success Messages
// =============================================================================

// Error messages.
const (
	// ErrToValidateRecord indicates a validation error.
	ErrToValidateRecord = "record validation error"

	// FailedToCreateRecord indicates failure to create a record.
	FailedToCreateRecord = "failed to create record"

	// FailedToGetRecord indicates failure to get a record.
	FailedToGetRecord = "failed to get record"

	// FailedToListRecords indicates failure to list records.
	FailedToListRecords = "failed to list records"

	// FailedToUpdateRecord indicates failure to update a record.
	FailedToUpdateRecord = "failed to update record"

	// FailedToDeleteRecord indicates failure to delete a record.
	FailedToDeleteRecord = "failed to delete record"

	// RecordNotFound indicates the record was not found.
	RecordNotFound = "record not found"
)

// Validation messages.
const (
	// UserIDIsRequired indicates the user ID is required.
	UserIDIsRequired = "user ID is required"

	// RecordIDIsRequired indicates the record ID is required.
	RecordIDIsRequired = "record ID is required"

	// TagIDIsRequired indicates the tag ID is required.
	TagIDIsRequired = "tag ID is required"

	// TagIDCannotBeZero indicates the tag ID cannot be zero.
	TagIDCannotBeZero = "tag id cannot be zero"

	// RecordedAtCannotBeInTheFuture indicates recordedAt must not be a future timestamp.
	RecordedAtCannotBeInTheFuture = "recordedAt cannot be in the future"

	// StartDateMustBeBeforeEndDate indicates start date validation error.
	StartDateMustBeBeforeEndDate = "startDate must be before or equal to endDate"

	// InvalidRecordIDOrUserID indicates invalid record or user ID.
	InvalidRecordIDOrUserID = "invalid recordID or userID"

	// UserNotAuthenticated indicates user is not authenticated.
	UserNotAuthenticated = "user not authenticated"

	// UserIDNegative indicates user ID cannot be negative.
	UserIDNegative = "user id negative"

	// UserIDStringNotSupported indicates string user IDs are not supported.
	UserIDStringNotSupported = "user id string not supported"

	// InvalidUserIDInContext indicates invalid user ID format in context.
	InvalidUserIDInContext = "invalid user id in context"
)

// =============================================================================
// DEFAULT VALUES - Business defaults for optional fields
// =============================================================================

const (
	// DefaultRecordStatus is the default status for new records.
	DefaultRecordStatus = "published"

	// DefaultTimezone is the default timezone when not provided.
	DefaultTimezone = "America/Sao_Paulo"
)

// =============================================================================
// SENTINEL ERRORS - For errors.Is() comparisons
// =============================================================================

var (
	// ErrValidateRecord is a sentinel error for record validation failures.
	ErrValidateRecord = errors.New(ErrToValidateRecord)

	// ErrCreateRecord is a sentinel error for record creation failures.
	ErrCreateRecord = errors.New(FailedToCreateRecord)

	// ErrGetRecord is a sentinel error for record retrieval failures.
	ErrGetRecord = errors.New(FailedToGetRecord)

	// ErrListRecords is a sentinel error for listing records.
	ErrListRecords = errors.New(FailedToListRecords)

	// ErrUpdateRecord is a sentinel error for record update failures.
	ErrUpdateRecord = errors.New(FailedToUpdateRecord)

	// ErrDeleteRecord is a sentinel error for record deletion failures.
	ErrDeleteRecord = errors.New(FailedToDeleteRecord)

	// ErrRecordNotFound is a sentinel error when record is not found.
	ErrRecordNotFound = errors.New(RecordNotFound)

	// ErrUserIDIsRequired is a sentinel error when user ID is missing.
	ErrUserIDIsRequired = errors.New(UserIDIsRequired)

	// ErrRecordIDIsRequired is a sentinel error when record ID is missing.
	ErrRecordIDIsRequired = errors.New(RecordIDIsRequired)

	// ErrTagIDIsRequired is a sentinel error when tag ID is missing.
	ErrTagIDIsRequired = errors.New(TagIDIsRequired)

	// ErrTagIDCannotBeZero is a sentinel error when tag ID is zero.
	ErrTagIDCannotBeZero = errors.New(TagIDCannotBeZero)

	// ErrRecordedAtFuture is a sentinel error when recordedAt is in future.
	ErrRecordedAtFuture = errors.New(RecordedAtCannotBeInTheFuture)

	// ErrStartDateAfterEndDate is a sentinel error for date range validation.
	ErrStartDateAfterEndDate = errors.New(StartDateMustBeBeforeEndDate)

	// ErrInvalidRecordIDOrUserID is a sentinel error for invalid IDs.
	ErrInvalidRecordIDOrUserID = errors.New(InvalidRecordIDOrUserID)

	// ErrUserNotAuthenticated is a sentinel error when user is not authenticated.
	ErrUserNotAuthenticated = errors.New(UserNotAuthenticated)

	// ErrUserIDNegative is a sentinel error when user ID is negative.
	ErrUserIDNegative = errors.New(UserIDNegative)

	// ErrUserIDStringNotSupported is a sentinel error for string user IDs.
	ErrUserIDStringNotSupported = errors.New(UserIDStringNotSupported)

	// ErrInvalidUserIDInContext is a sentinel error for invalid user ID in context.
	ErrInvalidUserIDInContext = errors.New(InvalidUserIDInContext)
)
