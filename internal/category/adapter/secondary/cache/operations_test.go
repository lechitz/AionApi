package cache_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	categorycache "github.com/lechitz/aion-api/internal/category/adapter/secondary/cache"
	"github.com/lechitz/aion-api/internal/category/core/domain"
	"github.com/stretchr/testify/require"
)

type fakeCategoryCache struct {
	data    map[string]string
	setErr  map[string]error
	getErr  map[string]error
	delErr  map[string]error
	lastTTL map[string]time.Duration
}

func newFakeCategoryCache() *fakeCategoryCache {
	return &fakeCategoryCache{
		data:    make(map[string]string),
		setErr:  make(map[string]error),
		getErr:  make(map[string]error),
		delErr:  make(map[string]error),
		lastTTL: make(map[string]time.Duration),
	}
}

func (f *fakeCategoryCache) Set(_ context.Context, key string, value interface{}, exp time.Duration) error {
	if err := f.setErr[key]; err != nil {
		return err
	}
	f.lastTTL[key] = exp
	f.data[key] = fmt.Sprintf("%v", value)
	return nil
}

func (f *fakeCategoryCache) Get(_ context.Context, key string) (string, error) {
	if err := f.getErr[key]; err != nil {
		return "", err
	}
	return f.data[key], nil
}

func (f *fakeCategoryCache) Del(_ context.Context, key string) error {
	if err := f.delErr[key]; err != nil {
		return err
	}
	delete(f.data, key)
	return nil
}

func (f *fakeCategoryCache) Ping(context.Context) error { return nil }
func (f *fakeCategoryCache) Close() error               { return nil }

type mockCategoryLogger struct{}

func (m *mockCategoryLogger) Infow(_ string, _ ...interface{})                        {}
func (m *mockCategoryLogger) Warnw(_ string, _ ...interface{})                        {}
func (m *mockCategoryLogger) Errorw(_ string, _ ...interface{})                       {}
func (m *mockCategoryLogger) InfowCtx(_ context.Context, _ string, _ ...interface{})  {}
func (m *mockCategoryLogger) WarnwCtx(_ context.Context, _ string, _ ...interface{})  {}
func (m *mockCategoryLogger) ErrorwCtx(_ context.Context, _ string, _ ...interface{}) {}
func (m *mockCategoryLogger) DebugwCtx(_ context.Context, _ string, _ ...interface{}) {}
func (m *mockCategoryLogger) Debugw(_ string, _ ...interface{})                       {}
func (m *mockCategoryLogger) Debugf(_ string, _ ...interface{})                       {}
func (m *mockCategoryLogger) Infof(_ string, _ ...interface{})                        {}
func (m *mockCategoryLogger) Warnf(_ string, _ ...interface{})                        {}
func (m *mockCategoryLogger) Errorf(_ string, _ ...interface{})                       {}
func (m *mockCategoryLogger) Sync() error                                             { return nil }

func testCategory() domain.Category {
	now := time.Date(2026, 2, 14, 12, 0, 0, 0, time.UTC)
	return domain.Category{
		ID:          10,
		UserID:      99,
		Name:        "Health",
		Description: "health routines",
		Color:       "#00AA77",
		Icon:        "heart",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestCategoryCache_SaveAndGetByIDAndName(t *testing.T) {
	fc := newFakeCategoryCache()
	store := categorycache.NewStore(fc, &mockCategoryLogger{})
	category := testCategory()

	err := store.SaveCategory(t.Context(), category, 0)
	require.NoError(t, err)
	idKey := fmt.Sprintf(categorycache.CategoryIDKeyFormat, category.UserID, category.ID)
	require.Equal(t, categorycache.CategoryExpirationDefault, fc.lastTTL[idKey])

	gotByID, err := store.GetCategory(t.Context(), category.ID, category.UserID)
	require.NoError(t, err)
	require.Equal(t, category.Name, gotByID.Name)

	err = store.SaveCategoryByName(t.Context(), category, 0)
	require.NoError(t, err)
	nameKey := fmt.Sprintf(categorycache.CategoryNameKeyFormat, category.UserID, category.Name)
	require.Equal(t, categorycache.CategoryExpirationDefault, fc.lastTTL[nameKey])

	gotByName, err := store.GetCategoryByName(t.Context(), category.Name, category.UserID)
	require.NoError(t, err)
	require.Equal(t, category.ID, gotByName.ID)
}

func TestCategoryCache_SaveAndGetList(t *testing.T) {
	fc := newFakeCategoryCache()
	store := categorycache.NewStore(fc, &mockCategoryLogger{})
	category := testCategory()

	require.NoError(t, store.SaveCategoryList(t.Context(), category.UserID, []domain.Category{category}, 0))
	listKey := fmt.Sprintf(categorycache.CategoryListKeyFormat, category.UserID)
	require.Equal(t, categorycache.CategoryListExpirationDefault, fc.lastTTL[listKey])

	got, err := store.GetCategoryList(t.Context(), category.UserID)
	require.NoError(t, err)
	require.Len(t, got, 1)
}

func TestCategoryCache_GetFallbackAndErrors(t *testing.T) {
	fc := newFakeCategoryCache()
	store := categorycache.NewStore(fc, &mockCategoryLogger{})

	emptyCategory, err := store.GetCategory(t.Context(), 1, 2)
	require.NoError(t, err)
	require.Equal(t, domain.Category{}, emptyCategory)

	emptyList, err := store.GetCategoryList(t.Context(), 2)
	require.NoError(t, err)
	require.Nil(t, emptyList)

	idKey := fmt.Sprintf(categorycache.CategoryIDKeyFormat, uint64(2), uint64(1))
	fc.getErr[idKey] = errors.New("redis get failed")
	_, err = store.GetCategory(t.Context(), 1, 2)
	require.Error(t, err)
	delete(fc.getErr, idKey)

	fc.data[idKey] = "{invalid-json"
	_, err = store.GetCategory(t.Context(), 1, 2)
	require.Error(t, err)

	listKey := fmt.Sprintf(categorycache.CategoryListKeyFormat, uint64(2))
	fc.data[listKey] = "not-json"
	_, err = store.GetCategoryList(t.Context(), 2)
	require.Error(t, err)
}

func TestCategoryCache_SaveAndDeleteErrors(t *testing.T) {
	fc := newFakeCategoryCache()
	store := categorycache.NewStore(fc, &mockCategoryLogger{})
	category := testCategory()

	idKey := fmt.Sprintf(categorycache.CategoryIDKeyFormat, category.UserID, category.ID)
	nameKey := fmt.Sprintf(categorycache.CategoryNameKeyFormat, category.UserID, category.Name)
	listKey := fmt.Sprintf(categorycache.CategoryListKeyFormat, category.UserID)

	fc.setErr[idKey] = errors.New("set id failed")
	require.Error(t, store.SaveCategory(t.Context(), category, time.Minute))
	delete(fc.setErr, idKey)

	fc.setErr[nameKey] = errors.New("set name failed")
	require.Error(t, store.SaveCategoryByName(t.Context(), category, time.Minute))
	delete(fc.setErr, nameKey)

	fc.setErr[listKey] = errors.New("set list failed")
	require.Error(t, store.SaveCategoryList(t.Context(), category.UserID, []domain.Category{category}, time.Minute))

	fc.delErr[idKey] = errors.New("del id failed")
	require.Error(t, store.DeleteCategory(t.Context(), category.ID, category.UserID))
	delete(fc.delErr, idKey)

	fc.delErr[nameKey] = errors.New("del name failed")
	require.Error(t, store.DeleteCategoryByName(t.Context(), category.Name, category.UserID))
	delete(fc.delErr, nameKey)

	fc.delErr[listKey] = errors.New("del list failed")
	require.Error(t, store.DeleteCategoryList(t.Context(), category.UserID))
}

func TestCategoryCache_DeleteSuccess(t *testing.T) {
	fc := newFakeCategoryCache()
	store := categorycache.NewStore(fc, &mockCategoryLogger{})
	category := testCategory()

	require.NoError(t, store.DeleteCategory(t.Context(), category.ID, category.UserID))
	require.NoError(t, store.DeleteCategoryByName(t.Context(), category.Name, category.UserID))
	require.NoError(t, store.DeleteCategoryList(t.Context(), category.UserID))
}
