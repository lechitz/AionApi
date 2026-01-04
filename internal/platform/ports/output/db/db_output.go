// Package db provides an abstraction for database operations.
package db

import (
	"context"
)

// DB is an abstraction for database operations following the Repository pattern.
// This interface allows the application layer to remain database-agnostic,
// following Hexagonal Architecture and Dependency Inversion principles.
//
// Implementations should provide:
// - Transaction management (Begin, Commit, Rollback)
// - Query execution (raw queries, ORM operations)
// - Connection lifecycle management.
type DB interface {
	// WithContext returns a new DB instance bound to the given context.
	// This is essential for request-scoped operations, timeouts, and cancellation.
	WithContext(ctx context.Context) DB

	// Create inserts a new record into the database.
	// The value parameter should be a pointer to the entity being created.
	Create(value interface{}) DB

	// Save updates all fields of an existing record or creates it if it doesn't exist.
	Save(value interface{}) DB

	// Updates specific fields of a record.
	// The updates parameter can be a map[string]interface{} or a struct with only the fields to update.
	Updates(updates interface{}) DB

	// Update updates a single column.
	Update(column string, value interface{}) DB

	// Delete performs a soft delete (if the model has a DeletedAt field) or hard delete.
	Delete(value interface{}) DB

	// First finds the first record ordered by primary key.
	First(dest interface{}, conds ...interface{}) DB

	// Find retrieves all records that match the given conditions.
	Find(dest interface{}, conds ...interface{}) DB

	// Model specifies the model to use for the query (for chaining with Updates, etc).
	Model(value interface{}) DB

	// Select specifies fields to retrieve from the database.
	Select(query interface{}, args ...interface{}) DB

	// Where adds a WHERE clause to the query.
	Where(query interface{}, args ...interface{}) DB

	// Order specifies the order when retrieving records.
	Order(value interface{}) DB

	// Limit specifies the maximum number of records to retrieve.
	Limit(limit int) DB

	// Offset specifies the number of records to skip before starting to return records.
	Offset(offset int) DB

	// Count counts the number of records that match the query.
	Count(count *int64) DB

	// Scan scans the result into the dest variable.
	Scan(dest interface{}) DB

	// Exec executes a raw SQL query (without returning rows).
	Exec(sql string, values ...interface{}) DB

	// Raw executes a raw SQL query and returns rows (for SELECT statements).
	// This allows complex queries like full-text search with ts_rank, JOINs, and GROUP BY.
	Raw(sql string, values ...interface{}) DB

	// Error returns the error from the last operation (if any).
	Error() error

	// RowsAffected returns the number of rows affected by the last operation.
	RowsAffected() int64

	// Begin starts a new transaction.
	Begin() DB

	// Commit commits the current transaction.
	Commit() DB

	// Rollback rolls back the current transaction.
	Rollback() DB

	// Transaction executes a function within a transaction.
	// If the function returns an error, the transaction is rolled back.
	// Otherwise, the transaction is committed.
	Transaction(fc func(tx DB) error) error
}
