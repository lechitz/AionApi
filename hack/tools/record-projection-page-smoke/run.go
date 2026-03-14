package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	loginPath   = "/aion/api/v1/auth/login"
	graphqlPath = "/aion/api/v1/graphql"
)

type loginEnvelope struct {
	Result struct {
		Token string `json:"token"`
	} `json:"result"`
}

type gqlEnvelope struct {
	Data   map[string]json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type createRecordResult struct {
	ID string `json:"id"`
}

type pageProjectionResult struct {
	RecordID      string  `json:"recordId"`
	Description   *string `json:"description"`
	EventTimeUTC  string  `json:"eventTimeUTC"`
	LastEventType string  `json:"lastEventType"`
}

type pageCursor struct {
	afterEventTime string
	afterID        string
}

type smokeRecord struct {
	id          string
	description string
	eventTime   string
}

const smokeDescriptionPrefix = "codex page smoke"

func run(ctx context.Context, cfg config) error {
	client := &http.Client{Timeout: cfg.timeout}

	token, err := login(ctx, client, cfg)
	if err != nil {
		return err
	}

	if err := cleanupExistingSmokeRecords(ctx, client, cfg, token); err != nil {
		return err
	}

	base := time.Now().UTC().Add(1 * time.Hour)
	records := []smokeRecord{
		{description: "codex page smoke oldest", eventTime: base.Format(time.RFC3339)},
		{description: "codex page smoke middle", eventTime: base.Add(1 * time.Minute).Format(time.RFC3339)},
		{description: "codex page smoke newest", eventTime: base.Add(2 * time.Minute).Format(time.RFC3339)},
	}

	for i := range records {
		recordID, createErr := createRecord(ctx, client, cfg, token, records[i].description, records[i].eventTime)
		if createErr != nil {
			return createErr
		}
		records[i].id = recordID
	}

	defer cleanupRecords(context.Background(), client, cfg, token, records)

	if err := waitUntilProjected(ctx, client, cfg, token, records); err != nil {
		return err
	}

	firstPage, err := fetchPage(ctx, client, cfg.host, token, cfg.pageLimit, nil)
	if err != nil {
		return err
	}
	if len(firstPage) != cfg.pageLimit {
		return fmt.Errorf("expected first page size %d, got %d", cfg.pageLimit, len(firstPage))
	}

	expectedFirstPage := []smokeRecord{records[2], records[1]}
	if err := assertPage(firstPage, expectedFirstPage); err != nil {
		return fmt.Errorf("first page validation failed: %w", err)
	}

	cursor := pageCursor{
		afterEventTime: firstPage[len(firstPage)-1].EventTimeUTC,
		afterID:        firstPage[len(firstPage)-1].RecordID,
	}
	secondPage, err := fetchPage(ctx, client, cfg.host, token, cfg.pageLimit, &cursor)
	if err != nil {
		return err
	}
	if len(secondPage) < 1 {
		return errors.New("expected at least one item on second page")
	}
	if err := assertPage(secondPage[:1], []smokeRecord{records[0]}); err != nil {
		return fmt.Errorf("second page leading item validation failed: %w", err)
	}

	if secondPage[0].RecordID == firstPage[0].RecordID || secondPage[0].RecordID == firstPage[1].RecordID {
		return errors.New("cursor pagination returned overlapping projection rows")
	}

	fmt.Printf("record projection page smoke passed for ids=%s,%s,%s\n", records[0].id, records[1].id, records[2].id)
	return nil
}

func cleanupExistingSmokeRecords(ctx context.Context, client *http.Client, cfg config, token string) error {
	page, err := fetchPage(ctx, client, cfg.host, token, 50, nil)
	if err != nil {
		return err
	}

	stale := make([]smokeRecord, 0)
	for _, item := range page {
		if item.Description == nil || !strings.HasPrefix(*item.Description, smokeDescriptionPrefix) {
			continue
		}
		stale = append(stale, smokeRecord{
			id:          item.RecordID,
			description: *item.Description,
			eventTime:   item.EventTimeUTC,
		})
	}

	if len(stale) == 0 {
		return nil
	}

	cleanupRecords(ctx, client, cfg, token, stale)
	return waitUntilRemoved(ctx, client, cfg, token, stale)
}

func login(ctx context.Context, client *http.Client, cfg config) (string, error) {
	body, err := json.Marshal(map[string]string{
		"username": cfg.username,
		"password": cfg.password,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(cfg.host, "/")+loginPath, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var envelope loginEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return "", err
	}
	if envelope.Result.Token == "" {
		return "", errors.New("login token missing")
	}
	return envelope.Result.Token, nil
}

func createRecord(ctx context.Context, client *http.Client, cfg config, token string, description string, eventTime string) (string, error) {
	query := fmt.Sprintf(`mutation { createRecord(input: { tagId: %q, description: %q, eventTime: %q, source: "codex-page-smoke", status: "published" }) { id } }`, cfg.tagID, description, eventTime)
	var result createRecordResult
	if err := graphql(ctx, client, cfg.host, token, query, "createRecord", &result); err != nil {
		return "", err
	}
	if result.ID == "" {
		return "", errors.New("createRecord returned empty id")
	}
	return result.ID, nil
}

func cleanupRecords(ctx context.Context, client *http.Client, cfg config, token string, records []smokeRecord) {
	for _, item := range records {
		if item.id == "" {
			continue
		}
		_ = deleteRecord(ctx, client, cfg, token, item.id)
	}
}

func waitUntilRemoved(ctx context.Context, client *http.Client, cfg config, token string, records []smokeRecord) error {
	deadline := time.Now().Add(cfg.timeout)
	for time.Now().Before(deadline) {
		allGone := true
		for _, item := range records {
			projection, err := fetchProjectionByID(ctx, client, cfg.host, token, item.id)
			if err != nil {
				allGone = false
				break
			}
			if projection != nil {
				allGone = false
				break
			}
		}
		if allGone {
			return nil
		}
		time.Sleep(cfg.pollInterval)
	}
	return errors.New("timed out waiting stale page smoke projections to be removed")
}

func deleteRecord(ctx context.Context, client *http.Client, cfg config, token string, recordID string) error {
	query := fmt.Sprintf(`mutation { softDeleteRecord(input: { id: %q }) }`, recordID)
	var deleted bool
	return graphql(ctx, client, cfg.host, token, query, "softDeleteRecord", &deleted)
}

func waitUntilProjected(ctx context.Context, client *http.Client, cfg config, token string, records []smokeRecord) error {
	deadline := time.Now().Add(cfg.timeout)
	for time.Now().Before(deadline) {
		if projected, err := allProjectedByID(ctx, client, cfg.host, token, records); err == nil && projected {
			return nil
		}
		time.Sleep(cfg.pollInterval)
	}
	return errors.New("timed out waiting for derived projections to contain all smoke records")
}

func allProjectedByID(ctx context.Context, client *http.Client, host string, token string, expected []smokeRecord) (bool, error) {
	for _, item := range expected {
		projection, err := fetchProjectionByID(ctx, client, host, token, item.id)
		if err != nil {
			return false, err
		}
		if projection == nil {
			return false, nil
		}
		if projection.Description == nil || *projection.Description != item.description {
			return false, nil
		}
		if projection.LastEventType != "record.created" {
			return false, nil
		}
	}

	return true, nil
}

func fetchProjectionByID(ctx context.Context, client *http.Client, host string, token string, recordID string) (*pageProjectionResult, error) {
	query := fmt.Sprintf(`query { recordProjectionById(id: %q) { recordId description eventTimeUTC lastEventType } }`, recordID)
	var result *pageProjectionResult
	if err := graphql(ctx, client, host, token, query, "recordProjectionById", &result); err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func fetchPage(ctx context.Context, client *http.Client, host string, token string, limit int, cursor *pageCursor) ([]pageProjectionResult, error) {
	query := `query ($limit: Int, $afterEventTime: String, $afterId: ID) {
		recordProjections(limit: $limit, afterEventTime: $afterEventTime, afterId: $afterId) {
			recordId
			description
			eventTimeUTC
			lastEventType
		}
	}`

	variables := map[string]any{
		"limit": limit,
	}
	if cursor != nil {
		variables["afterEventTime"] = cursor.afterEventTime
		variables["afterId"] = cursor.afterID
	}

	var result []pageProjectionResult
	if err := graphqlWithVariables(ctx, client, host, token, query, variables, "recordProjections", &result); err != nil {
		return nil, err
	}
	return result, nil
}

func assertPage(page []pageProjectionResult, expected []smokeRecord) error {
	if len(page) != len(expected) {
		return fmt.Errorf("expected %d items, got %d", len(expected), len(page))
	}

	for i := range expected {
		if page[i].RecordID != expected[i].id {
			return fmt.Errorf("position %d expected record_id=%s, got %s", i, expected[i].id, page[i].RecordID)
		}
		if page[i].Description == nil || *page[i].Description != expected[i].description {
			return fmt.Errorf("position %d expected description %q, got %v", i, expected[i].description, page[i].Description)
		}
		if page[i].LastEventType != "record.created" {
			return fmt.Errorf("position %d expected last_event_type record.created, got %s", i, page[i].LastEventType)
		}
	}

	return nil
}

func graphql(ctx context.Context, client *http.Client, host string, token string, query string, field string, out interface{}) error {
	return graphqlWithVariables(ctx, client, host, token, query, nil, field, out)
}

func graphqlWithVariables(ctx context.Context, client *http.Client, host string, token string, query string, variables map[string]any, field string, out interface{}) error {
	body, err := json.Marshal(map[string]any{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(host, "/")+graphqlPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var envelope gqlEnvelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return err
	}
	if len(envelope.Errors) > 0 {
		return errors.New(envelope.Errors[0].Message)
	}

	raw, ok := envelope.Data[field]
	if !ok {
		return fmt.Errorf("graphql field %s missing", field)
	}
	return json.Unmarshal(raw, out)
}
