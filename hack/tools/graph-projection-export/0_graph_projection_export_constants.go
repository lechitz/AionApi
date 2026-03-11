package main

import "time"

const (
	defaultExportWindow   = "WINDOW_30D"
	defaultExportTimezone = "America/Sao_Paulo"
	defaultExportLimit    = 5000

	seriesKeyRecordsCount = "records.count"
	jsonIndentPrefix      = ""
	jsonIndentValue       = "  "
	dateLayoutISO8601     = "2006-01-02"

	flagUserID     = "user-id"
	flagWindow     = "window"
	flagDate       = "date"
	flagTimezone   = "timezone"
	flagCategoryID = "category-id"
	flagTagIDs     = "tag-ids"
	flagOutput     = "output"

	logMsgExportStarted   = "graph projection export started"
	logMsgExportSucceeded = "graph projection export completed"
	logMsgExportFailed    = "graph projection export failed"
	logMsgOutputWritten   = "graph projection export written"

	logFieldComponent    = "component"
	logFieldWindow       = "window"
	logFieldTagIDsCount  = "tag_ids_count"
	logFieldNodeCount    = "node_count"
	logFieldEdgeCount    = "edge_count"
	logFieldOutput       = "output"
	logFieldGeneratedAt  = "generated_at"
	logFieldRecordsCount = "records_count"
	logFieldInsights     = "insights_count"
	logFieldDate         = "date"

	logValueComponent = "hack_graph_projection_export"
)

const (
	defaultDBConnectTimeout = 15 * time.Second
)
