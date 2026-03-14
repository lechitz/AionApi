// Package domain contains core business entities for the Chat context.
package domain

// ChatContext represents an aggregated view of user data for AI context.
// It provides recent activity and available resources to help the AI generate better responses.
type ChatContext struct {
	RecentRecords   []interface{} // Recent records created by the user
	Categories      []interface{} // Available categories for the user
	Tags            []interface{} // Available tags for the user
	RecentChats     []ChatHistory // Recent chat conversations
	TotalRecords    int           // Total number of records
	TotalCategories int           // Total number of categories
	TotalTags       int           // Total number of tags
}
