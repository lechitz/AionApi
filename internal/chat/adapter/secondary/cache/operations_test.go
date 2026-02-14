package cache_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	chatcache "github.com/lechitz/AionApi/internal/chat/adapter/secondary/cache"
	"github.com/lechitz/AionApi/internal/chat/core/domain"
	"github.com/stretchr/testify/require"
)

type fakeChatCache struct {
	data    map[string]string
	setErr  map[string]error
	getErr  map[string]error
	delErr  map[string]error
	lastTTL map[string]time.Duration
}

func newFakeChatCache() *fakeChatCache {
	return &fakeChatCache{
		data:    make(map[string]string),
		setErr:  make(map[string]error),
		getErr:  make(map[string]error),
		delErr:  make(map[string]error),
		lastTTL: make(map[string]time.Duration),
	}
}

func (f *fakeChatCache) Set(_ context.Context, key string, value interface{}, exp time.Duration) error {
	if err := f.setErr[key]; err != nil {
		return err
	}
	f.lastTTL[key] = exp
	f.data[key] = fmt.Sprintf("%v", value)
	return nil
}

func (f *fakeChatCache) Get(_ context.Context, key string) (string, error) {
	if err := f.getErr[key]; err != nil {
		return "", err
	}
	return f.data[key], nil
}

func (f *fakeChatCache) Del(_ context.Context, key string) error {
	if err := f.delErr[key]; err != nil {
		return err
	}
	delete(f.data, key)
	return nil
}

func (f *fakeChatCache) Ping(context.Context) error { return nil }
func (f *fakeChatCache) Close() error               { return nil }

type mockChatLogger struct{}

func (m *mockChatLogger) Infow(_ string, _ ...interface{})                        {}
func (m *mockChatLogger) Warnw(_ string, _ ...interface{})                        {}
func (m *mockChatLogger) Errorw(_ string, _ ...interface{})                       {}
func (m *mockChatLogger) InfowCtx(_ context.Context, _ string, _ ...interface{})  {}
func (m *mockChatLogger) WarnwCtx(_ context.Context, _ string, _ ...interface{})  {}
func (m *mockChatLogger) ErrorwCtx(_ context.Context, _ string, _ ...interface{}) {}
func (m *mockChatLogger) DebugwCtx(_ context.Context, _ string, _ ...interface{}) {}
func (m *mockChatLogger) Debugw(_ string, _ ...interface{})                       {}
func (m *mockChatLogger) Debugf(_ string, _ ...interface{})                       {}
func (m *mockChatLogger) Infof(_ string, _ ...interface{})                        {}
func (m *mockChatLogger) Warnf(_ string, _ ...interface{})                        {}
func (m *mockChatLogger) Errorf(_ string, _ ...interface{})                       {}
func (m *mockChatLogger) Sync() error                                             { return nil }

func chatEntry(idx uint64) domain.ChatHistory {
	now := time.Date(2026, 2, 14, 12, 0, 0, 0, time.UTC)
	return domain.ChatHistory{
		ChatID:     idx,
		UserID:     42,
		Message:    fmt.Sprintf("message-%d", idx),
		Response:   fmt.Sprintf("response-%d", idx),
		CreatedAt:  now,
		UpdatedAt:  now,
		TokensUsed: toIntSafe(idx),
	}
}

func toIntSafe(v uint64) int {
	if v > uint64(math.MaxInt) {
		return math.MaxInt
	}
	return int(v)
}

func TestChatCache_GetLatestMissAndError(t *testing.T) {
	fc := newFakeChatCache()
	store := chatcache.NewStore(fc, &mockChatLogger{})
	key := "chat:history:42"

	got, err := store.GetLatest(t.Context(), 42, 5)
	require.NoError(t, err)
	require.Empty(t, got)

	fc.getErr[key] = errors.New("redis down")
	_, err = store.GetLatest(t.Context(), 42, 5)
	require.Error(t, err)
}

func TestChatCache_GetLatestInvalidJSONAndLimit(t *testing.T) {
	fc := newFakeChatCache()
	store := chatcache.NewStore(fc, &mockChatLogger{})
	key := "chat:history:42"

	fc.data[key] = "{invalid-json"
	_, err := store.GetLatest(t.Context(), 42, 3)
	require.Error(t, err)

	entries := []domain.ChatHistory{chatEntry(1), chatEntry(2), chatEntry(3), chatEntry(4)}
	raw, err := json.Marshal(entries)
	require.NoError(t, err)
	fc.data[key] = string(raw)

	limited, err := store.GetLatest(t.Context(), 42, 2)
	require.NoError(t, err)
	require.Len(t, limited, 2)
	require.Equal(t, uint64(1), limited[0].ChatID)
	require.Equal(t, uint64(2), limited[1].ChatID)
}

func TestChatCache_AddPrependsAndTrims(t *testing.T) {
	fc := newFakeChatCache()
	store := chatcache.NewStore(fc, &mockChatLogger{})
	userID := uint64(42)
	key := "chat:history:42"

	for i := uint64(1); i <= 25; i++ {
		require.NoError(t, store.Add(t.Context(), userID, chatEntry(i)))
	}

	raw := fc.data[key]
	var stored []domain.ChatHistory
	require.NoError(t, json.Unmarshal([]byte(raw), &stored))
	require.Len(t, stored, 20)
	require.Equal(t, uint64(25), stored[0].ChatID)
	require.Equal(t, uint64(6), stored[19].ChatID)
}

func TestChatCache_AddFallbackAndSetError(t *testing.T) {
	fc := newFakeChatCache()
	store := chatcache.NewStore(fc, &mockChatLogger{})
	key := "chat:history:42"

	fc.getErr[key] = errors.New("read failed")
	err := store.Add(t.Context(), 42, chatEntry(1))
	require.NoError(t, err)

	fc.setErr[key] = errors.New("write failed")
	err = store.Add(t.Context(), 42, chatEntry(2))
	require.Error(t, err)
}

func TestChatCache_ClearSuccessAndError(t *testing.T) {
	fc := newFakeChatCache()
	store := chatcache.NewStore(fc, &mockChatLogger{})
	key := "chat:history:42"

	fc.data[key] = `[{"chatId":1,"userId":42,"message":"a","response":"b"}]`
	require.NoError(t, store.Clear(t.Context(), 42))
	require.Empty(t, fc.data[key])

	fc.delErr[key] = errors.New("del failed")
	err := store.Clear(t.Context(), 42)
	require.Error(t, err)
}
