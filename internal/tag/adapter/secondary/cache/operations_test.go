package cache_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	tagcache "github.com/lechitz/aion-api/internal/tag/adapter/secondary/cache"
	"github.com/lechitz/aion-api/internal/tag/core/domain"
	"github.com/stretchr/testify/require"
)

type fakeTagCache struct {
	data    map[string]string
	setErr  map[string]error
	getErr  map[string]error
	delErr  map[string]error
	lastTTL map[string]time.Duration
}

func newFakeTagCache() *fakeTagCache {
	return &fakeTagCache{
		data:    make(map[string]string),
		setErr:  make(map[string]error),
		getErr:  make(map[string]error),
		delErr:  make(map[string]error),
		lastTTL: make(map[string]time.Duration),
	}
}

func (f *fakeTagCache) Set(_ context.Context, key string, value interface{}, exp time.Duration) error {
	if err := f.setErr[key]; err != nil {
		return err
	}
	f.lastTTL[key] = exp
	f.data[key] = fmt.Sprintf("%v", value)
	return nil
}

func (f *fakeTagCache) Get(_ context.Context, key string) (string, error) {
	if err := f.getErr[key]; err != nil {
		return "", err
	}
	return f.data[key], nil
}

func (f *fakeTagCache) Del(_ context.Context, key string) error {
	if err := f.delErr[key]; err != nil {
		return err
	}
	delete(f.data, key)
	return nil
}

func (f *fakeTagCache) Ping(context.Context) error { return nil }
func (f *fakeTagCache) Close() error               { return nil }

type mockTagLogger struct{}

func (m *mockTagLogger) Infow(_ string, _ ...interface{})                        {}
func (m *mockTagLogger) Warnw(_ string, _ ...interface{})                        {}
func (m *mockTagLogger) Errorw(_ string, _ ...interface{})                       {}
func (m *mockTagLogger) InfowCtx(_ context.Context, _ string, _ ...interface{})  {}
func (m *mockTagLogger) WarnwCtx(_ context.Context, _ string, _ ...interface{})  {}
func (m *mockTagLogger) ErrorwCtx(_ context.Context, _ string, _ ...interface{}) {}
func (m *mockTagLogger) DebugwCtx(_ context.Context, _ string, _ ...interface{}) {}
func (m *mockTagLogger) Debugw(_ string, _ ...interface{})                       {}
func (m *mockTagLogger) Debugf(_ string, _ ...interface{})                       {}
func (m *mockTagLogger) Infof(_ string, _ ...interface{})                        {}
func (m *mockTagLogger) Warnf(_ string, _ ...interface{})                        {}
func (m *mockTagLogger) Errorf(_ string, _ ...interface{})                       {}
func (m *mockTagLogger) Sync() error                                             { return nil }

func testTag() domain.Tag {
	now := time.Date(2026, 2, 14, 12, 0, 0, 0, time.UTC)
	return domain.Tag{
		ID:          10,
		UserID:      99,
		CategoryID:  7,
		Name:        "Health",
		Description: "health routines",
		Icon:        "H",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestTagCache_SaveAndGetByIDAndName(t *testing.T) {
	fc := newFakeTagCache()
	store := tagcache.NewStore(fc, &mockTagLogger{})
	tag := testTag()

	err := store.SaveTag(t.Context(), tag, 0)
	require.NoError(t, err)
	idKey := fmt.Sprintf(tagcache.TagIDKeyFormat, tag.UserID, tag.ID)
	require.Equal(t, tagcache.TagExpirationDefault, fc.lastTTL[idKey])

	gotByID, err := store.GetTag(t.Context(), tag.ID, tag.UserID)
	require.NoError(t, err)
	require.Equal(t, tag.Name, gotByID.Name)

	err = store.SaveTagByName(t.Context(), tag, 0)
	require.NoError(t, err)
	nameKey := fmt.Sprintf(tagcache.TagNameKeyFormat, tag.UserID, tag.Name)
	require.Equal(t, tagcache.TagExpirationDefault, fc.lastTTL[nameKey])

	gotByName, err := store.GetTagByName(t.Context(), tag.Name, tag.UserID)
	require.NoError(t, err)
	require.Equal(t, tag.ID, gotByName.ID)
}

func TestTagCache_SaveAndGetLists(t *testing.T) {
	fc := newFakeTagCache()
	store := tagcache.NewStore(fc, &mockTagLogger{})
	tag := testTag()
	tags := []domain.Tag{tag}

	require.NoError(t, store.SaveTagList(t.Context(), tag.UserID, tags, 0))
	listKey := fmt.Sprintf(tagcache.TagListKeyFormat, tag.UserID)
	require.Equal(t, tagcache.TagListExpirationDefault, fc.lastTTL[listKey])

	gotList, err := store.GetTagList(t.Context(), tag.UserID)
	require.NoError(t, err)
	require.Len(t, gotList, 1)

	require.NoError(t, store.SaveTagsByCategory(t.Context(), tag.CategoryID, tag.UserID, tags, 0))
	catKey := fmt.Sprintf(tagcache.TagByCategoryKeyFormat, tag.CategoryID, tag.UserID)
	require.Equal(t, tagcache.TagListExpirationDefault, fc.lastTTL[catKey])

	gotByCategory, err := store.GetTagsByCategory(t.Context(), tag.CategoryID, tag.UserID)
	require.NoError(t, err)
	require.Len(t, gotByCategory, 1)
}

func TestTagCache_GetFallbackAndErrors(t *testing.T) {
	fc := newFakeTagCache()
	store := tagcache.NewStore(fc, &mockTagLogger{})

	emptyTag, err := store.GetTag(t.Context(), 1, 2)
	require.NoError(t, err)
	require.Equal(t, domain.Tag{}, emptyTag)

	emptyList, err := store.GetTagList(t.Context(), 2)
	require.NoError(t, err)
	require.Nil(t, emptyList)

	idKey := fmt.Sprintf(tagcache.TagIDKeyFormat, uint64(2), uint64(1))
	fc.getErr[idKey] = errors.New("redis get failed")
	_, err = store.GetTag(t.Context(), 1, 2)
	require.Error(t, err)
	delete(fc.getErr, idKey)

	fc.data[idKey] = "{invalid-json"
	_, err = store.GetTag(t.Context(), 1, 2)
	require.Error(t, err)

	listKey := fmt.Sprintf(tagcache.TagListKeyFormat, uint64(2))
	fc.data[listKey] = "not-json"
	_, err = store.GetTagList(t.Context(), 2)
	require.Error(t, err)
}

func TestTagCache_GetByNameAndCategoryErrors(t *testing.T) {
	fc := newFakeTagCache()
	store := tagcache.NewStore(fc, &mockTagLogger{})

	nameKey := fmt.Sprintf(tagcache.TagNameKeyFormat, uint64(2), "Health")
	fc.getErr[nameKey] = errors.New("redis get name failed")
	_, err := store.GetTagByName(t.Context(), "Health", 2)
	require.Error(t, err)

	delete(fc.getErr, nameKey)
	fc.data[nameKey] = "{invalid-json"
	_, err = store.GetTagByName(t.Context(), "Health", 2)
	require.Error(t, err)

	catKey := fmt.Sprintf(tagcache.TagByCategoryKeyFormat, uint64(7), uint64(2))
	fc.getErr[catKey] = errors.New("redis get category failed")
	_, err = store.GetTagsByCategory(t.Context(), 7, 2)
	require.Error(t, err)

	delete(fc.getErr, catKey)
	fc.data[catKey] = "not-json"
	_, err = store.GetTagsByCategory(t.Context(), 7, 2)
	require.Error(t, err)
}

func TestTagCache_SaveAndDeleteErrors(t *testing.T) {
	fc := newFakeTagCache()
	store := tagcache.NewStore(fc, &mockTagLogger{})
	tag := testTag()

	idKey := fmt.Sprintf(tagcache.TagIDKeyFormat, tag.UserID, tag.ID)
	nameKey := fmt.Sprintf(tagcache.TagNameKeyFormat, tag.UserID, tag.Name)
	listKey := fmt.Sprintf(tagcache.TagListKeyFormat, tag.UserID)
	catKey := fmt.Sprintf(tagcache.TagByCategoryKeyFormat, tag.CategoryID, tag.UserID)

	fc.setErr[idKey] = errors.New("set id failed")
	require.Error(t, store.SaveTag(t.Context(), tag, time.Minute))
	delete(fc.setErr, idKey)

	fc.setErr[nameKey] = errors.New("set name failed")
	require.Error(t, store.SaveTagByName(t.Context(), tag, time.Minute))
	delete(fc.setErr, nameKey)

	fc.setErr[listKey] = errors.New("set list failed")
	require.Error(t, store.SaveTagList(t.Context(), tag.UserID, []domain.Tag{tag}, time.Minute))
	delete(fc.setErr, listKey)

	fc.setErr[catKey] = errors.New("set category failed")
	require.Error(t, store.SaveTagsByCategory(t.Context(), tag.CategoryID, tag.UserID, []domain.Tag{tag}, time.Minute))

	fc.delErr[idKey] = errors.New("del id failed")
	require.Error(t, store.DeleteTag(t.Context(), tag.ID, tag.UserID))
	delete(fc.delErr, idKey)

	fc.delErr[nameKey] = errors.New("del name failed")
	require.Error(t, store.DeleteTagByName(t.Context(), tag.Name, tag.UserID))
	delete(fc.delErr, nameKey)

	fc.delErr[listKey] = errors.New("del list failed")
	require.Error(t, store.DeleteTagList(t.Context(), tag.UserID))
	delete(fc.delErr, listKey)

	fc.delErr[catKey] = errors.New("del category failed")
	require.Error(t, store.DeleteTagsByCategory(t.Context(), tag.CategoryID, tag.UserID))
}

func TestTagCache_DeleteSuccess(t *testing.T) {
	fc := newFakeTagCache()
	store := tagcache.NewStore(fc, &mockTagLogger{})
	tag := testTag()

	require.NoError(t, store.DeleteTag(t.Context(), tag.ID, tag.UserID))
	require.NoError(t, store.DeleteTagByName(t.Context(), tag.Name, tag.UserID))
	require.NoError(t, store.DeleteTagList(t.Context(), tag.UserID))
	require.NoError(t, store.DeleteTagsByCategory(t.Context(), tag.CategoryID, tag.UserID))
}
