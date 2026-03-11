package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	recorddomain "github.com/lechitz/AionApi/internal/record/core/domain"
)

var (
	errUserIDRequired    = errors.New("user id is required")
	errInvalidDate       = errors.New("invalid date")
	errInvalidCategoryID = errors.New("invalid category id")
	errInvalidTagIDs     = errors.New("invalid tag ids")
	errInvalidWindow     = errors.New("invalid window")
)

type exportConfig struct {
	UserID     uint64
	Window     recorddomain.InsightWindow
	Date       time.Time
	Timezone   string
	CategoryID *uint64
	TagIDs     []uint64
	OutputPath string
}

func loadConfig(args []string) (exportConfig, error) {
	fs := flag.NewFlagSet("graph-projection-export", flag.ContinueOnError)

	var (
		userIDRaw     uint64
		windowRaw     string
		dateRaw       string
		timezoneRaw   string
		categoryIDRaw string
		tagIDsRaw     string
		outputPathRaw string
	)

	fs.Uint64Var(&userIDRaw, flagUserID, 0, "user id to export")
	fs.StringVar(&windowRaw, flagWindow, defaultExportWindow, "analysis window: WINDOW_7D, WINDOW_30D or WINDOW_90D")
	fs.StringVar(&dateRaw, flagDate, "", "target date in YYYY-MM-DD; defaults to today in timezone")
	fs.StringVar(&timezoneRaw, flagTimezone, defaultExportTimezone, "IANA timezone name")
	fs.StringVar(&categoryIDRaw, flagCategoryID, "", "optional category id scope")
	fs.StringVar(&tagIDsRaw, flagTagIDs, "", "optional comma-separated tag ids scope")
	fs.StringVar(&outputPathRaw, flagOutput, "", "optional output file path; defaults to stdout")

	if err := fs.Parse(args); err != nil {
		return exportConfig{}, err
	}

	if userIDRaw == 0 {
		return exportConfig{}, errUserIDRequired
	}

	window, err := parseWindow(windowRaw)
	if err != nil {
		return exportConfig{}, err
	}

	timezone := strings.TrimSpace(timezoneRaw)
	if timezone == "" {
		timezone = defaultExportTimezone
	}

	dateValue, err := parseDate(dateRaw, timezone)
	if err != nil {
		return exportConfig{}, err
	}

	categoryID, err := parseOptionalUint64(categoryIDRaw, errInvalidCategoryID)
	if err != nil {
		return exportConfig{}, err
	}

	tagIDs, err := parseTagIDs(tagIDsRaw)
	if err != nil {
		return exportConfig{}, err
	}

	return exportConfig{
		UserID:     userIDRaw,
		Window:     window,
		Date:       dateValue,
		Timezone:   timezone,
		CategoryID: categoryID,
		TagIDs:     tagIDs,
		OutputPath: strings.TrimSpace(outputPathRaw),
	}, nil
}

func parseWindow(raw string) (recorddomain.InsightWindow, error) {
	switch strings.TrimSpace(strings.ToUpper(raw)) {
	case string(recorddomain.InsightWindow7D):
		return recorddomain.InsightWindow7D, nil
	case string(recorddomain.InsightWindow30D):
		return recorddomain.InsightWindow30D, nil
	case string(recorddomain.InsightWindow90D):
		return recorddomain.InsightWindow90D, nil
	default:
		return "", fmt.Errorf("%w: %s", errInvalidWindow, raw)
	}
}

func parseDate(raw string, timezone string) (time.Time, error) {
	location, err := time.LoadLocation(strings.TrimSpace(timezone))
	if err != nil {
		location = time.UTC
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		now := time.Now().In(location)
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location), nil
	}

	parsed, err := time.ParseInLocation(dateLayoutISO8601, raw, location)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s", errInvalidDate, raw)
	}

	return parsed, nil
}

func parseOptionalUint64(raw string, sentinel error) (*uint64, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil || parsed == 0 {
		return nil, fmt.Errorf("%w: %s", sentinel, raw)
	}
	return &parsed, nil
}

func parseTagIDs(raw string) ([]uint64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	tagIDs := make([]uint64, 0, len(parts))
	seen := make(map[uint64]struct{}, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil || parsed == 0 {
			return nil, fmt.Errorf("%w: %s", errInvalidTagIDs, value)
		}
		if _, exists := seen[parsed]; exists {
			continue
		}
		seen[parsed] = struct{}{}
		tagIDs = append(tagIDs, parsed)
	}
	return tagIDs, nil
}
