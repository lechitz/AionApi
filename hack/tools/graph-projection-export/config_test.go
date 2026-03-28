package main

import (
	"errors"
	"testing"

	recorddomain "github.com/lechitz/aion-api/internal/record/core/domain"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		wantErr error
		assert  func(t *testing.T, got exportConfig)
	}{
		{
			name:    "returns error when user id is missing",
			args:    []string{},
			wantErr: errUserIDRequired,
		},
		{
			name: "loads defaults with required user id",
			args: []string{"--user-id", "999"},
			assert: func(t *testing.T, got exportConfig) {
				t.Helper()
				if got.UserID != 999 {
					t.Fatalf("unexpected user id: %d", got.UserID)
				}
				if got.Window != recorddomain.InsightWindow30D {
					t.Fatalf("unexpected window: %s", got.Window)
				}
				if got.Timezone != defaultExportTimezone {
					t.Fatalf("unexpected timezone: %s", got.Timezone)
				}
			},
		},
		{
			name: "loads scoped filters",
			args: []string{
				"--user-id", "999",
				"--window", "WINDOW_7D",
				"--date", "2026-03-10",
				"--timezone", "UTC",
				"--category-id", "10",
				"--tag-ids", "2,4,4",
				"--output", "tmp/export.json",
			},
			assert: func(t *testing.T, got exportConfig) {
				t.Helper()
				if got.Window != recorddomain.InsightWindow7D {
					t.Fatalf("unexpected window: %s", got.Window)
				}
				if got.CategoryID == nil || *got.CategoryID != 10 {
					t.Fatalf("unexpected category id: %#v", got.CategoryID)
				}
				if len(got.TagIDs) != 2 || got.TagIDs[0] != 2 || got.TagIDs[1] != 4 {
					t.Fatalf("unexpected tag ids: %#v", got.TagIDs)
				}
				if got.OutputPath != "tmp/export.json" {
					t.Fatalf("unexpected output path: %s", got.OutputPath)
				}
			},
		},
		{
			name:    "returns error on invalid tag ids",
			args:    []string{"--user-id", "999", "--tag-ids", "abc"},
			wantErr: errInvalidTagIDs,
		},
		{
			name:    "returns error on invalid window",
			args:    []string{"--user-id", "999", "--window", "LAST_5D"},
			wantErr: errInvalidWindow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := loadConfig(tt.args)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.assert != nil {
				tt.assert(t, got)
			}
		})
	}
}
