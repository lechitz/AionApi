package cache_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	recordcache "github.com/lechitz/aion-api/internal/record/adapter/secondary/cache"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/stretchr/testify/require"
)

type fakeRecordCache struct {
	data    map[string]string
	setErr  map[string]error
	getErr  map[string]error
	delErr  map[string]error
	lastTTL map[string]time.Duration
}

func newFakeRecordCache() *fakeRecordCache {
	return &fakeRecordCache{
		data:    make(map[string]string),
		setErr:  make(map[string]error),
		getErr:  make(map[string]error),
		delErr:  make(map[string]error),
		lastTTL: make(map[string]time.Duration),
	}
}

func (f *fakeRecordCache) Set(_ context.Context, key string, value interface{}, exp time.Duration) error {
	if err := f.setErr[key]; err != nil {
		return err
	}
	f.lastTTL[key] = exp
	f.data[key] = fmt.Sprintf("%v", value)
	return nil
}

func (f *fakeRecordCache) Get(_ context.Context, key string) (string, error) {
	if err := f.getErr[key]; err != nil {
		return "", err
	}
	return f.data[key], nil
}

func (f *fakeRecordCache) Del(_ context.Context, key string) error {
	if err := f.delErr[key]; err != nil {
		return err
	}
	delete(f.data, key)
	return nil
}

func (f *fakeRecordCache) Ping(context.Context) error { return nil }
func (f *fakeRecordCache) Close() error               { return nil }

type mockRecordLogger struct{}

func (m *mockRecordLogger) Infow(_ string, _ ...interface{})                        {}
func (m *mockRecordLogger) Warnw(_ string, _ ...interface{})                        {}
func (m *mockRecordLogger) Errorw(_ string, _ ...interface{})                       {}
func (m *mockRecordLogger) InfowCtx(_ context.Context, _ string, _ ...interface{})  {}
func (m *mockRecordLogger) WarnwCtx(_ context.Context, _ string, _ ...interface{})  {}
func (m *mockRecordLogger) ErrorwCtx(_ context.Context, _ string, _ ...interface{}) {}
func (m *mockRecordLogger) DebugwCtx(_ context.Context, _ string, _ ...interface{}) {}
func (m *mockRecordLogger) Debugw(_ string, _ ...interface{})                       {}
func (m *mockRecordLogger) Debugf(_ string, _ ...interface{})                       {}
func (m *mockRecordLogger) Infof(_ string, _ ...interface{})                        {}
func (m *mockRecordLogger) Warnf(_ string, _ ...interface{})                        {}
func (m *mockRecordLogger) Errorf(_ string, _ ...interface{})                       {}
func (m *mockRecordLogger) Sync() error                                             { return nil }

func testRecord() domain.Record {
	desc := "gym"
	value := 10.5
	now := time.Date(2026, 2, 14, 12, 0, 0, 0, time.UTC)
	return domain.Record{
		ID:          1,
		UserID:      99,
		TagID:       7,
		Description: &desc,
		EventTime:   now,
		Value:       &value,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestRecordCache_SaveAndGetRecord(t *testing.T) {
	fc := newFakeRecordCache()
	store := recordcache.NewStore(fc, &mockRecordLogger{})
	rec := testRecord()

	err := store.SaveRecord(t.Context(), rec, 0)
	require.NoError(t, err)

	key := fmt.Sprintf(recordcache.RecordIDKeyFormat, rec.UserID, rec.ID)
	require.Equal(t, recordcache.RecordExpirationDefault, fc.lastTTL[key])

	got, err := store.GetRecord(t.Context(), rec.ID, rec.UserID)
	require.NoError(t, err)
	require.Equal(t, rec.ID, got.ID)
	require.Equal(t, rec.UserID, got.UserID)
}

func TestRecordCache_SaveAndGetLists(t *testing.T) {
	fc := newFakeRecordCache()
	store := recordcache.NewStore(fc, &mockRecordLogger{})
	rec := testRecord()
	records := []domain.Record{rec}
	date := rec.EventTime

	require.NoError(t, store.SaveRecordsByDay(t.Context(), rec.UserID, date, records, 0))
	dayKey := fmt.Sprintf(recordcache.RecordDayKeyFormat, rec.UserID, date.Format("2006-01-02"))
	require.Equal(t, recordcache.RecordListExpirationDefault, fc.lastTTL[dayKey])

	gotDay, err := store.GetRecordsByDay(t.Context(), rec.UserID, date)
	require.NoError(t, err)
	require.Len(t, gotDay, 1)

	require.NoError(t, store.SaveRecordsByCategory(t.Context(), 12, rec.UserID, records, 0))
	gotCategory, err := store.GetRecordsByCategory(t.Context(), 12, rec.UserID)
	require.NoError(t, err)
	require.Len(t, gotCategory, 1)

	require.NoError(t, store.SaveRecordsByTag(t.Context(), rec.TagID, rec.UserID, records, 0))
	gotTag, err := store.GetRecordsByTag(t.Context(), rec.TagID, rec.UserID)
	require.NoError(t, err)
	require.Len(t, gotTag, 1)
}

func TestRecordCache_GetFallbackAndErrors(t *testing.T) {
	fc := newFakeRecordCache()
	store := recordcache.NewStore(fc, &mockRecordLogger{})

	empty, err := store.GetRecord(t.Context(), 1, 2)
	require.NoError(t, err)
	require.Equal(t, domain.Record{}, empty)

	emptyList, err := store.GetRecordsByDay(t.Context(), 2, time.Now().UTC())
	require.NoError(t, err)
	require.Nil(t, emptyList)

	key := fmt.Sprintf(recordcache.RecordIDKeyFormat, uint64(2), uint64(1))
	fc.getErr[key] = errors.New("redis get failed")
	_, err = store.GetRecord(t.Context(), 1, 2)
	require.Error(t, err)

	fc.data[key] = "{invalid-json"
	delete(fc.getErr, key)
	_, err = store.GetRecord(t.Context(), 1, 2)
	require.Error(t, err)
}

func TestRecordCache_GetCategoryAndTagErrors(t *testing.T) {
	fc := newFakeRecordCache()
	store := recordcache.NewStore(fc, &mockRecordLogger{})

	categoryKey := fmt.Sprintf(recordcache.RecordByCategoryKeyFormat, uint64(5), uint64(3))
	fc.getErr[categoryKey] = errors.New("redis get category failed")
	_, err := store.GetRecordsByCategory(t.Context(), 5, 3)
	require.Error(t, err)

	delete(fc.getErr, categoryKey)
	fc.data[categoryKey] = "not-json"
	_, err = store.GetRecordsByCategory(t.Context(), 5, 3)
	require.Error(t, err)

	tagKey := fmt.Sprintf(recordcache.RecordByTagKeyFormat, uint64(7), uint64(3))
	fc.getErr[tagKey] = errors.New("redis get tag failed")
	_, err = store.GetRecordsByTag(t.Context(), 7, 3)
	require.Error(t, err)

	delete(fc.getErr, tagKey)
	fc.data[tagKey] = "not-json"
	_, err = store.GetRecordsByTag(t.Context(), 7, 3)
	require.Error(t, err)
}

func TestRecordCache_SaveSetError(t *testing.T) {
	fc := newFakeRecordCache()
	store := recordcache.NewStore(fc, &mockRecordLogger{})
	rec := testRecord()
	key := fmt.Sprintf(recordcache.RecordIDKeyFormat, rec.UserID, rec.ID)
	fc.setErr[key] = errors.New("redis set failed")
	err := store.SaveRecord(t.Context(), rec, time.Minute)
	require.Error(t, err)
}

func TestRecordCache_ListSetGetDeleteErrors(t *testing.T) {
	fc := newFakeRecordCache()
	store := recordcache.NewStore(fc, &mockRecordLogger{})
	date := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)

	dayKey := fmt.Sprintf(recordcache.RecordDayKeyFormat, uint64(3), date.Format("2006-01-02"))
	fc.setErr[dayKey] = errors.New("set day failed")
	err := store.SaveRecordsByDay(t.Context(), 3, date, []domain.Record{testRecord()}, time.Minute)
	require.Error(t, err)
	delete(fc.setErr, dayKey)

	fc.data[dayKey] = "not-json"
	_, err = store.GetRecordsByDay(t.Context(), 3, date)
	require.Error(t, err)

	categoryKey := fmt.Sprintf(recordcache.RecordByCategoryKeyFormat, uint64(5), uint64(3))
	fc.setErr[categoryKey] = errors.New("set category failed")
	err = store.SaveRecordsByCategory(t.Context(), 5, 3, []domain.Record{testRecord()}, time.Minute)
	require.Error(t, err)

	tagKey := fmt.Sprintf(recordcache.RecordByTagKeyFormat, uint64(7), uint64(3))
	fc.setErr[tagKey] = errors.New("set tag failed")
	err = store.SaveRecordsByTag(t.Context(), 7, 3, []domain.Record{testRecord()}, time.Minute)
	require.Error(t, err)

	fc.delErr[dayKey] = errors.New("delete day failed")
	err = store.DeleteRecordsByDay(t.Context(), 3, date)
	require.Error(t, err)
}

func TestRecordCache_DeleteAllKeys(t *testing.T) {
	fc := newFakeRecordCache()
	store := recordcache.NewStore(fc, &mockRecordLogger{})
	date := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)

	require.NoError(t, store.DeleteRecord(t.Context(), 10, 20))
	require.NoError(t, store.DeleteRecordsByDay(t.Context(), 20, date))
	require.NoError(t, store.DeleteRecordsByCategory(t.Context(), 30, 20))
	require.NoError(t, store.DeleteRecordsByTag(t.Context(), 40, 20))

	rec := testRecord()
	raw, err := json.Marshal([]domain.Record{rec})
	require.NoError(t, err)
	key := fmt.Sprintf(recordcache.RecordByTagKeyFormat, rec.TagID, rec.UserID)
	fc.data[key] = string(raw)
	_, err = store.GetRecordsByTag(t.Context(), rec.TagID, rec.UserID)
	require.NoError(t, err)
}

func TestRecordCache_DeleteErrorsByKeyType(t *testing.T) {
	fc := newFakeRecordCache()
	store := recordcache.NewStore(fc, &mockRecordLogger{})

	recordKey := fmt.Sprintf(recordcache.RecordIDKeyFormat, uint64(20), uint64(10))
	fc.delErr[recordKey] = errors.New("delete record failed")
	err := store.DeleteRecord(t.Context(), 10, 20)
	require.Error(t, err)

	categoryKey := fmt.Sprintf(recordcache.RecordByCategoryKeyFormat, uint64(30), uint64(20))
	fc.delErr[categoryKey] = errors.New("delete category list failed")
	err = store.DeleteRecordsByCategory(t.Context(), 30, 20)
	require.Error(t, err)

	tagKey := fmt.Sprintf(recordcache.RecordByTagKeyFormat, uint64(40), uint64(20))
	fc.delErr[tagKey] = errors.New("delete tag list failed")
	err = store.DeleteRecordsByTag(t.Context(), 40, 20)
	require.Error(t, err)
}
