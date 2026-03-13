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

type projectionResult struct {
	RecordID      string  `json:"recordId"`
	Description   *string `json:"description"`
	LastEventType string  `json:"lastEventType"`
}

func run(ctx context.Context, cfg config) error {
	client := &http.Client{Timeout: cfg.timeout}

	token, err := login(ctx, client, cfg)
	if err != nil {
		return err
	}

	recordID, err := createRecord(ctx, client, cfg, token)
	if err != nil {
		return err
	}

	if err := waitProjection(ctx, client, cfg, token, recordID, "record.created", "codex smoke created", true); err != nil {
		return err
	}
	if err := updateRecord(ctx, client, cfg, token, recordID); err != nil {
		return err
	}
	if err := waitProjection(ctx, client, cfg, token, recordID, "record.updated", "codex smoke updated", true); err != nil {
		return err
	}
	if err := deleteRecord(ctx, client, cfg, token, recordID); err != nil {
		return err
	}
	if err := waitProjection(ctx, client, cfg, token, recordID, "record.deleted", "", false); err != nil {
		return err
	}

	fmt.Printf("record projection smoke passed for record_id=%s\n", recordID)
	return nil
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

func createRecord(ctx context.Context, client *http.Client, cfg config, token string) (string, error) {
	query := fmt.Sprintf(`mutation { createRecord(input: { tagId: %q, description: "codex smoke created", source: "codex-smoke", status: "published" }) { id } }`, cfg.tagID)
	var result createRecordResult
	if err := graphql(ctx, client, cfg.host, token, query, "createRecord", &result); err != nil {
		return "", err
	}
	if result.ID == "" {
		return "", errors.New("createRecord returned empty id")
	}
	return result.ID, nil
}

func updateRecord(ctx context.Context, client *http.Client, cfg config, token string, recordID string) error {
	query := fmt.Sprintf(`mutation { updateRecord(input: { id: %q, description: "codex smoke updated", source: "codex-smoke-updated" }) { id } }`, recordID)
	var result createRecordResult
	return graphql(ctx, client, cfg.host, token, query, "updateRecord", &result)
}

func deleteRecord(ctx context.Context, client *http.Client, cfg config, token string, recordID string) error {
	query := fmt.Sprintf(`mutation { softDeleteRecord(input: { id: %q }) }`, recordID)
	var deleted bool
	return graphql(ctx, client, cfg.host, token, query, "softDeleteRecord", &deleted)
}

func waitProjection(ctx context.Context, client *http.Client, cfg config, token string, recordID string, expectedEventType string, expectedDescription string, expectPresent bool) error {
	deadline := time.Now().Add(cfg.timeout)
	for time.Now().Before(deadline) {
		projection, err := fetchProjection(ctx, client, cfg.host, token, recordID)
		if err == nil {
			if expectPresent && projection != nil && projection.LastEventType == expectedEventType && projection.Description != nil && *projection.Description == expectedDescription {
				return nil
			}
			if !expectPresent && projection == nil {
				return nil
			}
		}
		time.Sleep(cfg.pollInterval)
	}

	if expectPresent {
		return fmt.Errorf("projection for record_id=%s did not reach expected state %s", recordID, expectedEventType)
	}
	return fmt.Errorf("projection for record_id=%s still present after delete", recordID)
}

func fetchProjection(ctx context.Context, client *http.Client, host string, token string, recordID string) (*projectionResult, error) {
	query := fmt.Sprintf(`query { recordProjectionById(id: %q) { recordId description lastEventType } }`, recordID)
	var result *projectionResult
	if err := graphql(ctx, client, host, token, query, "recordProjectionById", &result); err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func graphql(ctx context.Context, client *http.Client, host string, token string, query string, field string, out interface{}) error {
	body, err := json.Marshal(map[string]string{"query": query})
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
