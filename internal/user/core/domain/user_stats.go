// Package domain contains core business entities for the User context.
package domain

// UserStats represents aggregated statistics about a user's data.
type UserStats struct {
	TotalRecords     int         // Total number of records
	TotalCategories  int         // Total number of categories
	TotalTags        int         // Total number of tags
	RecordsThisWeek  int         // Records created this week
	RecordsThisMonth int         // Records created this month
	MostUsedCategory *UsageCount // Most frequently used category
	MostUsedTag      *UsageCount // Most frequently used tag
}

// UsageCount represents a generic entity usage count.
// This type is domain-agnostic and can represent categories, tags, or other entities
// without coupling the User bounded context to other domains.
type UsageCount struct {
	ID    uint64 // Entity ID
	Name  string // Entity name
	Count int    // Usage count
}
