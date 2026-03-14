package cache_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/user/adapter/secondary/cache"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserCacheDTO_NoPasswordHash verifies that UserCacheDTO never includes PasswordHash field.
// SECURITY TEST: Critical - ensures password hashes are never cached.
func TestUserCacheDTO_NoPasswordHash(t *testing.T) {
	dto := cache.UserCacheDTO{
		Version:   cache.UserCacheSchemaVersion,
		ID:        123,
		Name:      "Test User",
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Serialize DTO to JSON (what gets stored in Redis)
	data, err := json.Marshal(dto)
	require.NoError(t, err, "DTO serialization should succeed")

	// Parse JSON to map to inspect fields
	var fields map[string]interface{}
	err = json.Unmarshal(data, &fields)
	require.NoError(t, err, "JSON parsing should succeed")

	// SECURITY: Verify password_hash field does NOT exist
	_, hasPasswordHash := fields["password_hash"]
	assert.False(t, hasPasswordHash, "SECURITY VIOLATION: password_hash found in cached data!")

	_, hasPassword := fields["password"]
	assert.False(t, hasPassword, "SECURITY VIOLATION: password found in cached data!")

	// Verify expected fields ARE present
	assert.InDelta(t, 123, fields["id"], 0.001, "ID should be present")
	assert.Equal(t, "testuser", fields["username"], "Username should be present")
	assert.Equal(t, "test@example.com", fields["email"], "Email should be present")
}

// TestSaveUser_StripsPasswordHash verifies that SaveUser never caches password hash.
// SECURITY TEST: Critical - ensures domain.User with PasswordHash doesn't leak to cache.
func TestSaveUser_StripsPasswordHash(t *testing.T) {
	// Create a mock cache that captures what gets stored
	mockCache := &mockCacheCapture{data: make(map[string]string)}
	mockLogger := &mockLogger{}

	store := cache.NewStore(mockCache, mockLogger)

	// Create domain.User (domain.User no longer has PasswordHash field - it's in a separate table)
	user := domain.User{
		ID:        456,
		Name:      "Secure User",
		Username:  "secureuser",
		Email:     "secure@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to cache
	ctx := t.Context()
	err := store.SaveUser(ctx, user, time.Minute)
	require.NoError(t, err, "SaveUser should succeed")

	// Verify what was actually stored in cache
	cachedData := mockCache.data["user:id:456"]
	require.NotEmpty(t, cachedData, "Data should be cached")

	// Parse cached JSON
	var cachedFields map[string]interface{}
	err = json.Unmarshal([]byte(cachedData), &cachedFields)
	require.NoError(t, err, "Cached data should be valid JSON")

	// SECURITY: Verify password hash was NOT cached
	_, hasPasswordHash := cachedFields["password_hash"]
	assert.False(t, hasPasswordHash, "SECURITY VIOLATION: password_hash was cached!")

	_, hasPassword := cachedFields["password"]
	assert.False(t, hasPassword, "SECURITY VIOLATION: password was cached!")

	// Verify safe fields ARE cached
	assert.InDelta(t, 456, cachedFields["id"], 0.001)
	assert.Equal(t, "secureuser", cachedFields["username"])
	assert.Equal(t, "secure@example.com", cachedFields["email"])
}

// TestGetUser_NeverHasPasswordHash verifies that retrieved user never has password hash.
// SECURITY TEST: Ensures cached users never have password hash populated.
func TestGetUser_NeverHasPasswordHash(t *testing.T) {
	// Prepare cache with user data (without password)
	mockCache := &mockCacheCapture{data: make(map[string]string)}
	mockLogger := &mockLogger{}

	userData := `{
		"version": 2,
		"id": 789,
		"name": "Cached User",
		"username": "cacheduser",
		"email": "cached@example.com",
		"created_at": "2025-01-22T00:00:00Z",
		"updated_at": "2025-01-22T00:00:00Z"
	}`
	mockCache.data["user:id:789"] = userData

	store := cache.NewStore(mockCache, mockLogger)

	// Get user from cache
	ctx := t.Context()
	user, err := store.GetUserByID(ctx, 789)
	require.NoError(t, err, "GetUserByID should succeed")

	// SECURITY: domain.User no longer has PasswordHash field - passwords stored separately
	// Verify other fields are populated
	assert.Equal(t, uint64(789), user.ID)
	assert.Equal(t, "cacheduser", user.Username)
	assert.Equal(t, "cached@example.com", user.Email)
}

func TestGetUserByID_LegacyPayloadReturnsError(t *testing.T) {
	mockCache := &mockCacheCapture{data: make(map[string]string)}
	mockLogger := &mockLogger{}

	// Legacy payload without "version" should be treated as stale.
	mockCache.data["user:id:999"] = `{
		"id": 999,
		"name": "Legacy User",
		"username": "legacy",
		"email": "legacy@example.com"
	}`

	store := cache.NewStore(mockCache, mockLogger)

	_, err := store.GetUserByID(t.Context(), 999)
	require.Error(t, err, "legacy cache payload must trigger fallback to repository")
}

// mockCacheCapture captures what gets stored in cache for inspection.
type mockCacheCapture struct {
	data map[string]string
}

func (m *mockCacheCapture) Set(_ context.Context, key string, value interface{}, _ time.Duration) error {
	// Convert value to string for storage in our mock
	var strValue string
	if str, ok := value.(string); ok {
		strValue = str
	} else {
		// For testing purposes, we expect string values
		strValue = fmt.Sprintf("%v", value)
	}
	m.data[key] = strValue
	return nil
}

func (m *mockCacheCapture) Get(_ context.Context, key string) (string, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return "", assert.AnError
}

func (m *mockCacheCapture) Del(_ context.Context, key string) error {
	delete(m.data, key)
	return nil
}

func (m *mockCacheCapture) Ping(_ context.Context) error {
	return nil
}

func (m *mockCacheCapture) Close() error {
	return nil
}

// mockLogger prevents actual logging during tests.
type mockLogger struct{}

func (m *mockLogger) Infow(_ string, _ ...interface{})                        {}
func (m *mockLogger) Warnw(_ string, _ ...interface{})                        {}
func (m *mockLogger) Errorw(_ string, _ ...interface{})                       {}
func (m *mockLogger) InfowCtx(_ context.Context, _ string, _ ...interface{})  {}
func (m *mockLogger) WarnwCtx(_ context.Context, _ string, _ ...interface{})  {}
func (m *mockLogger) ErrorwCtx(_ context.Context, _ string, _ ...interface{}) {}
func (m *mockLogger) DebugwCtx(_ context.Context, _ string, _ ...interface{}) {}
func (m *mockLogger) Debugw(_ string, _ ...interface{})                       {}
func (m *mockLogger) Debugf(_ string, _ ...interface{})                       {}
func (m *mockLogger) Infof(_ string, _ ...interface{})                        {}
func (m *mockLogger) Warnf(_ string, _ ...interface{})                        {}
func (m *mockLogger) Errorf(_ string, _ ...interface{})                       {}
func (m *mockLogger) Sync() error                                             { return nil }
