package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
)

const (
	realtimeLoginPath   = "/aion/api/v1/auth/login"
	realtimeGraphqlPath = "/aion/api/v1/graphql"
	realtimeStreamPath  = "/aion/api/v1/realtime/events/stream"
)

var errRealtimeProjectionNotFound = errors.New("record projection not found")

type realtimeLoginEnvelope struct {
	Result struct {
		Token string `json:"token"`
	} `json:"result"`
}

type realtimeGQLEnvelope struct {
	Data   map[string]json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type realtimeCreateRecordResult struct {
	ID string `json:"id"`
}

type realtimeProjectionResult struct {
	RecordID      string  `json:"recordId"`
	Description   *string `json:"description"`
	LastEventType string  `json:"lastEventType"`
}

type sseEnvelope struct {
	Event string
	Data  string
}

type projectionChangedPayload struct {
	Type           string      `json:"type"`
	UserID         interface{} `json:"userId"`
	RecordID       interface{} `json:"recordId"`
	Action         string      `json:"action"`
	ProjectedAtUTC string      `json:"projectedAtUTC"`
	SourceEventID  string      `json:"sourceEventId"`
	TraceID        string      `json:"traceId"`
	RequestID      string      `json:"requestId"`
}

func run(ctx context.Context, cfg config) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: cfg.timeout, Jar: jar}
	streamClient := &http.Client{Jar: jar}

	token, err := realtimeLogin(ctx, client, cfg)
	if err != nil {
		return err
	}

	streamCtx, cancel := context.WithTimeout(ctx, cfg.timeout)
	defer cancel()

	events := make(chan sseEnvelope, 8)
	errs := make(chan error, 1)
	go consumeSSE(streamCtx, streamClient, cfg, events, errs)

	if err := waitConnected(streamCtx, events, errs); err != nil {
		return err
	}

	recordID, err := realtimeCreateRecord(ctx, client, cfg, token)
	if err != nil {
		return err
	}
	defer func() {
		_ = realtimeDeleteRecord(context.Background(), client, cfg, token, recordID)
	}()

	payload, err := waitProjectionChanged(streamCtx, events, errs, recordID)
	if err != nil {
		return err
	}
	if payload.Action != "created" {
		return fmt.Errorf("unexpected realtime action %q for record_id=%s", payload.Action, recordID)
	}

	if err := waitRealtimeProjection(ctx, client, cfg, token, recordID); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(os.Stdout, "realtime record smoke passed for record_id=%s action=%s\n", recordID, payload.Action); err != nil {
		return err
	}
	return nil
}

func realtimeLogin(ctx context.Context, client *http.Client, cfg config) (string, error) {
	body, err := json.Marshal(map[string]string{
		"username": cfg.username,
		"password": cfg.password,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(cfg.host, "/")+realtimeLoginPath, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var envelope realtimeLoginEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return "", err
	}
	if envelope.Result.Token == "" {
		return "", errors.New("login token missing")
	}
	return envelope.Result.Token, nil
}

func consumeSSE(ctx context.Context, client *http.Client, cfg config, events chan<- sseEnvelope, errs chan<- error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(cfg.host, "/")+realtimeStreamPath, nil)
	if err != nil {
		errs <- err
		return
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := client.Do(req)
	if err != nil {
		errs <- err
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		errs <- fmt.Errorf("sse stream returned status=%d body=%s", resp.StatusCode, string(body))
		return
	}

	reader := bufio.NewScanner(resp.Body)
	var currentEvent string
	var dataLines []string

	for reader.Scan() {
		line := reader.Text()
		if line == "" {
			if currentEvent != "" {
				events <- sseEnvelope{
					Event: currentEvent,
					Data:  strings.Join(dataLines, "\n"),
				}
			}
			currentEvent = ""
			dataLines = dataLines[:0]
			continue
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		if strings.HasPrefix(line, "event:") {
			currentEvent = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			continue
		}
		if strings.HasPrefix(line, "data:") {
			dataLines = append(dataLines, strings.TrimSpace(strings.TrimPrefix(line, "data:")))
		}
	}

	if err := reader.Err(); err != nil && !errors.Is(err, context.Canceled) {
		errs <- err
	}
}

func waitConnected(ctx context.Context, events <-chan sseEnvelope, errs <-chan error) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting connected event: %w", ctx.Err())
		case err := <-errs:
			return err
		case event := <-events:
			if event.Event == "connected" {
				return nil
			}
		}
	}
}

func realtimeCreateRecord(ctx context.Context, client *http.Client, cfg config, token string) (string, error) {
	query := fmt.Sprintf(
		`mutation { createRecord(input: { tagId: %q, description: "codex realtime smoke", source: "codex-realtime-smoke", status: "published" }) { id } }`,
		cfg.tagID,
	)
	var result realtimeCreateRecordResult
	if err := realtimeGraphql(ctx, client, cfg.host, token, query, "createRecord", &result); err != nil {
		return "", err
	}
	if result.ID == "" {
		return "", errors.New("createRecord returned empty id")
	}
	return result.ID, nil
}

func realtimeDeleteRecord(ctx context.Context, client *http.Client, cfg config, token string, recordID string) error {
	query := fmt.Sprintf(`mutation { softDeleteRecord(input: { id: %q }) }`, recordID)
	var deleted bool
	return realtimeGraphql(ctx, client, cfg.host, token, query, "softDeleteRecord", &deleted)
}

func waitProjectionChanged(ctx context.Context, events <-chan sseEnvelope, errs <-chan error, recordID string) (projectionChangedPayload, error) {
	for {
		select {
		case <-ctx.Done():
			return projectionChangedPayload{}, fmt.Errorf("waiting record_projection_changed: %w", ctx.Err())
		case err := <-errs:
			return projectionChangedPayload{}, err
		case event := <-events:
			if event.Event != "record_projection_changed" {
				continue
			}
			var payload projectionChangedPayload
			if err := json.Unmarshal([]byte(event.Data), &payload); err != nil {
				return projectionChangedPayload{}, err
			}
			if fmt.Sprint(payload.RecordID) == recordID {
				return payload, nil
			}
		}
	}
}

func waitRealtimeProjection(ctx context.Context, client *http.Client, cfg config, token string, recordID string) error {
	deadline := time.Now().Add(cfg.timeout)
	for time.Now().Before(deadline) {
		projection, err := fetchRealtimeProjection(ctx, client, cfg.host, token, recordID)
		if errors.Is(err, errRealtimeProjectionNotFound) {
			time.Sleep(cfg.pollInterval)
			continue
		}
		if err == nil && projection != nil && projection.LastEventType == "record.created" && projection.Description != nil &&
			*projection.Description == "codex realtime smoke" {
			return nil
		}
		time.Sleep(cfg.pollInterval)
	}
	return fmt.Errorf("projection for record_id=%s did not reach expected realtime state", recordID)
}

func fetchRealtimeProjection(ctx context.Context, client *http.Client, host string, token string, recordID string) (*realtimeProjectionResult, error) {
	query := fmt.Sprintf(`query { recordProjectionById(id: %q) { recordId description lastEventType } }`, recordID)
	var result *realtimeProjectionResult
	if err := realtimeGraphql(ctx, client, host, token, query, "recordProjectionById", &result); err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, errRealtimeProjectionNotFound
		}
		return nil, err
	}
	return result, nil
}

func realtimeGraphql(ctx context.Context, client *http.Client, host string, token string, query string, field string, out interface{}) error {
	body, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(host, "/")+realtimeGraphqlPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var envelope realtimeGQLEnvelope
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
