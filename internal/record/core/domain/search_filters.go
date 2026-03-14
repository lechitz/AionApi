// Package domain contains core business entities for the Record context.
package domain

import "time"

// SearchFilters represents filters for searching records.
type SearchFilters struct {
	Query       string     // Search query (full-text search)
	CategoryIDs []uint64   // Filter by category IDs
	TagIDs      []uint64   // Filter by tag IDs
	StartDate   *time.Time // Filter records from this date
	EndDate     *time.Time // Filter records until this date
	Limit       int        // Maximum number of results
	Offset      int        // Number of results to skip (pagination)
}
