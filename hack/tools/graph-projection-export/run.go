package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lechitz/AionApi/internal/adapter/secondary/contextlogger"
	"github.com/lechitz/AionApi/internal/adapter/secondary/crypto"
	"github.com/lechitz/AionApi/internal/adapter/secondary/db/postgres"
	categoryRepo "github.com/lechitz/AionApi/internal/category/adapter/secondary/db/repository"
	categorydomain "github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/platform/config"
	recordRepo "github.com/lechitz/AionApi/internal/record/adapter/secondary/db/repository"
	recorddomain "github.com/lechitz/AionApi/internal/record/core/domain"
	recordinput "github.com/lechitz/AionApi/internal/record/core/ports/input"
	record "github.com/lechitz/AionApi/internal/record/core/usecase"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	tagRepo "github.com/lechitz/AionApi/internal/tag/adapter/secondary/db/repository"
	tagdomain "github.com/lechitz/AionApi/internal/tag/core/domain"
)

func run(ctx context.Context, args []string) int {
	cfg, err := loadConfig(args)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		return 1
	}

	log, cleanup := contextlogger.New()
	_ = cleanup

	appCfg, err := config.New(crypto.New()).Load(log)
	if err != nil {
		log.Errorw(logMsgExportFailed, commonkeys.Error, err.Error(), logFieldComponent, logValueComponent)
		return 1
	}
	if err := appCfg.Validate(); err != nil {
		log.Errorw(logMsgExportFailed, commonkeys.Error, err.Error(), logFieldComponent, logValueComponent)
		return 1
	}

	dbCtx, cancel := context.WithTimeout(ctx, defaultDBConnectTimeout)
	defer cancel()

	gormDB, err := postgres.NewConnection(dbCtx, appCfg.DB, log)
	if err != nil {
		log.Errorw(logMsgExportFailed, commonkeys.Error, err.Error(), logFieldComponent, logValueComponent)
		return 1
	}
	defer postgres.Close(gormDB, log)

	dbAdapter := postgres.NewDBAdapter(gormDB)
	categoryRepository := categoryRepo.New(dbAdapter, log)
	tagRepository := tagRepo.New(dbAdapter, log)
	recordRepository := recordRepo.New(dbAdapter, log)
	recordService := record.NewService(recordRepository, nil, tagRepository, log)

	log.Infow(
		logMsgExportStarted,
		logFieldComponent, logValueComponent,
		commonkeys.UserID, cfg.UserID,
		logFieldWindow, cfg.Window,
		logFieldDate, cfg.Date.Format(dateLayoutISO8601),
		commonkeys.Timezone, cfg.Timezone,
		logFieldTagIDsCount, len(cfg.TagIDs),
	)

	projection, err := buildProjection(ctx, cfg, categoryRepository, tagRepository, recordRepository, recordService)
	if err != nil {
		log.Errorw(
			logMsgExportFailed,
			logFieldComponent, logValueComponent,
			commonkeys.UserID, cfg.UserID,
			commonkeys.Error, err.Error(),
		)
		return 1
	}

	payload, err := json.MarshalIndent(projection, jsonIndentPrefix, jsonIndentValue)
	if err != nil {
		log.Errorw(logMsgExportFailed, commonkeys.Error, err.Error(), logFieldComponent, logValueComponent)
		return 1
	}

	if err := writeOutput(cfg.OutputPath, payload); err != nil {
		log.Errorw(logMsgExportFailed, commonkeys.Error, err.Error(), logFieldComponent, logValueComponent)
		return 1
	}

	log.Infow(
		logMsgExportSucceeded,
		logFieldComponent, logValueComponent,
		commonkeys.UserID, cfg.UserID,
		logFieldWindow, cfg.Window,
		logFieldNodeCount, len(projection.Nodes),
		logFieldEdgeCount, len(projection.Edges),
		logFieldGeneratedAt, projection.GeneratedAt.Format(time.RFC3339),
	)
	if strings.TrimSpace(cfg.OutputPath) != "" {
		log.Infow(logMsgOutputWritten, logFieldOutput, cfg.OutputPath, logFieldComponent, logValueComponent)
	}
	return 0
}

func buildProjection(
	ctx context.Context,
	cfg exportConfig,
	categoryRepository interface {
		ListAll(context.Context, uint64) ([]categorydomain.Category, error)
	},
	tagRepository interface {
		GetAll(context.Context, uint64) ([]tagdomain.Tag, error)
	},
	recordRepository interface {
		ListAllBetween(context.Context, uint64, time.Time, time.Time, int) ([]recorddomain.Record, error)
	},
	recordService interface {
		InsightFeed(context.Context, uint64, recordinput.InsightFeedQuery) ([]recorddomain.InsightCard, error)
	},
) (recorddomain.GraphProjection, error) {
	startUTC, endUTC := windowRange(cfg.Date, cfg.Timezone, cfg.Window)

	categories, err := categoryRepository.ListAll(ctx, cfg.UserID)
	if err != nil {
		return recorddomain.GraphProjection{}, err
	}
	tags, err := tagRepository.GetAll(ctx, cfg.UserID)
	if err != nil {
		return recorddomain.GraphProjection{}, err
	}
	records, err := recordRepository.ListAllBetween(ctx, cfg.UserID, startUTC, endUTC, defaultExportLimit)
	if err != nil {
		return recorddomain.GraphProjection{}, err
	}

	filteredRecords := filterRecordsByScope(records, cfg.CategoryID, cfg.TagIDs, tags)
	filteredTags := filterTagsByScope(tags, filteredRecords, cfg.CategoryID, cfg.TagIDs)
	filteredCategories := filterCategoriesByScope(categories, filteredTags, cfg.CategoryID)

	insights, err := recordService.InsightFeed(ctx, cfg.UserID, recordinput.InsightFeedQuery{
		Window:     string(cfg.Window),
		Limit:      defaultExportLimit,
		Date:       cfg.Date,
		Timezone:   cfg.Timezone,
		CategoryID: cfg.CategoryID,
		TagIDs:     cfg.TagIDs,
	})
	if err != nil {
		return recorddomain.GraphProjection{}, err
	}

	supportedRecordIDs := make(map[string][]uint64, len(insights))
	scopedTagIDs := make(map[string][]uint64, len(insights))
	scopedCategoryIDs := make(map[string][]uint64, len(insights))
	recordIDs := collectRecordIDs(filteredRecords)

	for _, insight := range insights {
		if len(recordIDs) > 0 {
			supportedRecordIDs[insight.ID] = append([]uint64(nil), recordIDs...)
		}
		if len(cfg.TagIDs) > 0 {
			scopedTagIDs[insight.ID] = append([]uint64(nil), cfg.TagIDs...)
		}
		if cfg.CategoryID != nil {
			scopedCategoryIDs[insight.ID] = []uint64{*cfg.CategoryID}
		}
	}

	return record.BuildGraphProjection(record.GraphProjectionBuildInput{
		UserID:                    cfg.UserID,
		GeneratedAt:               time.Now().UTC(),
		Timezone:                  cfg.Timezone,
		Records:                   filteredRecords,
		Tags:                      filteredTags,
		Categories:                filteredCategories,
		Insights:                  insights,
		InsightSupportedRecordIDs: supportedRecordIDs,
		InsightScopedTagIDs:       scopedTagIDs,
		InsightScopedCategoryIDs:  scopedCategoryIDs,
	}), nil
}

func windowRange(targetDate time.Time, timezone string, window recorddomain.InsightWindow) (time.Time, time.Time) {
	location, err := time.LoadLocation(strings.TrimSpace(timezone))
	if err != nil || location == nil {
		location = time.UTC
	}
	localDate := targetDate.In(location)
	startLocal := localDate.AddDate(0, 0, -(windowDays(window) - 1))
	endLocal := time.Date(localDate.Year(), localDate.Month(), localDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), location)
	return time.Date(startLocal.Year(), startLocal.Month(), startLocal.Day(), 0, 0, 0, 0, location).UTC(), endLocal.UTC()
}

func windowDays(window recorddomain.InsightWindow) int {
	switch window {
	case recorddomain.InsightWindow30D:
		return 30
	case recorddomain.InsightWindow90D:
		return 90
	default:
		return 7
	}
}

func filterRecordsByScope(records []recorddomain.Record, categoryID *uint64, tagIDs []uint64, tags []tagdomain.Tag) []recorddomain.Record {
	if categoryID == nil && len(tagIDs) == 0 {
		return records
	}

	tagCategoryByID := make(map[uint64]uint64, len(tags))
	for _, tag := range tags {
		tagCategoryByID[tag.ID] = tag.CategoryID
	}

	tagFilter := make(map[uint64]struct{}, len(tagIDs))
	for _, tagID := range tagIDs {
		tagFilter[tagID] = struct{}{}
	}

	filtered := make([]recorddomain.Record, 0, len(records))
	for _, recordItem := range records {
		if len(tagFilter) > 0 {
			if _, ok := tagFilter[recordItem.TagID]; !ok {
				continue
			}
		}
		if categoryID != nil {
			if tagCategoryByID[recordItem.TagID] != *categoryID {
				continue
			}
		}
		filtered = append(filtered, recordItem)
	}
	return filtered
}

func filterTagsByScope(tags []tagdomain.Tag, records []recorddomain.Record, categoryID *uint64, tagIDs []uint64) []tagdomain.Tag {
	if categoryID == nil && len(tagIDs) == 0 && len(records) == 0 {
		return tags
	}

	allowedByRecord := make(map[uint64]struct{}, len(records))
	for _, recordItem := range records {
		allowedByRecord[recordItem.TagID] = struct{}{}
	}
	allowedByFilter := make(map[uint64]struct{}, len(tagIDs))
	for _, tagID := range tagIDs {
		allowedByFilter[tagID] = struct{}{}
	}

	filtered := make([]tagdomain.Tag, 0, len(tags))
	for _, tag := range tags {
		if categoryID != nil && tag.CategoryID != *categoryID {
			continue
		}
		if len(allowedByFilter) > 0 {
			if _, ok := allowedByFilter[tag.ID]; !ok {
				continue
			}
		} else if len(allowedByRecord) > 0 {
			if _, ok := allowedByRecord[tag.ID]; !ok {
				continue
			}
		}
		filtered = append(filtered, tag)
	}
	return filtered
}

func filterCategoriesByScope(categories []categorydomain.Category, tags []tagdomain.Tag, categoryID *uint64) []categorydomain.Category {
	if categoryID != nil {
		filtered := make([]categorydomain.Category, 0, 1)
		for _, category := range categories {
			if category.ID == *categoryID {
				filtered = append(filtered, category)
				return filtered
			}
		}
		return filtered
	}

	if len(tags) == 0 {
		return categories
	}

	allowed := make(map[uint64]struct{}, len(tags))
	for _, tag := range tags {
		allowed[tag.CategoryID] = struct{}{}
	}

	filtered := make([]categorydomain.Category, 0, len(allowed))
	for _, category := range categories {
		if _, ok := allowed[category.ID]; ok {
			filtered = append(filtered, category)
		}
	}
	return filtered
}

func collectRecordIDs(records []recorddomain.Record) []uint64 {
	recordIDs := make([]uint64, 0, len(records))
	for _, recordItem := range records {
		recordIDs = append(recordIDs, recordItem.ID)
	}
	return recordIDs
}

func writeOutput(outputPath string, payload []byte) error {
	if strings.TrimSpace(outputPath) == "" {
		_, err := os.Stdout.Write(append(payload, '\n'))
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outputPath, append(payload, '\n'), 0o644)
}
