package controller_test

import (
	"context"
	"errors"
	"testing"
	"time"

	gmodel "github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/tag/adapter/primary/graphql/controller"
	"github.com/lechitz/aion-api/internal/tag/core/domain"
	taginput "github.com/lechitz/aion-api/internal/tag/core/ports/input"
	"github.com/stretchr/testify/require"
)

type tagServiceStub struct {
	createFn        func(context.Context, taginput.CreateTagCommand) (domain.Tag, error)
	updateFn        func(context.Context, taginput.UpdateTagCommand) (domain.Tag, error)
	getByIDFn       func(context.Context, uint64, uint64) (domain.Tag, error)
	getByNameFn     func(context.Context, string, uint64) (domain.Tag, error)
	getByCategoryFn func(context.Context, uint64, uint64) ([]domain.Tag, error)
	getAllFn        func(context.Context, uint64) ([]domain.Tag, error)
	softDeleteFn    func(context.Context, uint64, uint64) error
}

func (s *tagServiceStub) Create(ctx context.Context, cmd taginput.CreateTagCommand) (domain.Tag, error) {
	if s.createFn == nil {
		panic("unexpected Create call")
	}
	return s.createFn(ctx, cmd)
}

func (s *tagServiceStub) Update(ctx context.Context, cmd taginput.UpdateTagCommand) (domain.Tag, error) {
	if s.updateFn == nil {
		panic("unexpected Update call")
	}
	return s.updateFn(ctx, cmd)
}

func (s *tagServiceStub) GetByID(ctx context.Context, tagID uint64, userID uint64) (domain.Tag, error) {
	if s.getByIDFn == nil {
		panic("unexpected GetByID call")
	}
	return s.getByIDFn(ctx, tagID, userID)
}

func (s *tagServiceStub) GetByName(ctx context.Context, tagName string, userID uint64) (domain.Tag, error) {
	if s.getByNameFn == nil {
		panic("unexpected GetByName call")
	}
	return s.getByNameFn(ctx, tagName, userID)
}

func (s *tagServiceStub) GetByCategoryID(ctx context.Context, categoryID uint64, userID uint64) ([]domain.Tag, error) {
	if s.getByCategoryFn == nil {
		panic("unexpected GetByCategoryID call")
	}
	return s.getByCategoryFn(ctx, categoryID, userID)
}

func (s *tagServiceStub) GetAll(ctx context.Context, userID uint64) ([]domain.Tag, error) {
	if s.getAllFn == nil {
		panic("unexpected GetAll call")
	}
	return s.getAllFn(ctx, userID)
}

func (s *tagServiceStub) SoftDelete(ctx context.Context, tagID, userID uint64) error {
	if s.softDeleteFn == nil {
		panic("unexpected SoftDelete call")
	}
	return s.softDeleteFn(ctx, tagID, userID)
}

type tagLoggerStub struct{}

func (tagLoggerStub) Infof(string, ...any)                      {}
func (tagLoggerStub) Errorf(string, ...any)                     {}
func (tagLoggerStub) Debugf(string, ...any)                     {}
func (tagLoggerStub) Warnf(string, ...any)                      {}
func (tagLoggerStub) Infow(string, ...any)                      {}
func (tagLoggerStub) Errorw(string, ...any)                     {}
func (tagLoggerStub) Debugw(string, ...any)                     {}
func (tagLoggerStub) Warnw(string, ...any)                      {}
func (tagLoggerStub) InfowCtx(context.Context, string, ...any)  {}
func (tagLoggerStub) ErrorwCtx(context.Context, string, ...any) {}
func (tagLoggerStub) WarnwCtx(context.Context, string, ...any)  {}
func (tagLoggerStub) DebugwCtx(context.Context, string, ...any) {}

func TestTagController_Basics(t *testing.T) {
	now := time.Date(2026, 2, 14, 12, 0, 0, 0, time.UTC)
	desc := "d"
	icon := "i"
	newName := "Updated"
	newCategoryID := "2"

	ctrl := controller.NewController(&tagServiceStub{
		createFn: func(_ context.Context, cmd taginput.CreateTagCommand) (domain.Tag, error) {
			require.Equal(t, uint64(9), cmd.UserID)
			require.Equal(t, uint64(2), cmd.CategoryID)
			return domain.Tag{ID: 1, UserID: cmd.UserID, CategoryID: cmd.CategoryID, Name: cmd.Name, Description: desc, Icon: icon, CreatedAt: now, UpdatedAt: now}, nil
		},
		updateFn: func(_ context.Context, cmd taginput.UpdateTagCommand) (domain.Tag, error) {
			require.Equal(t, uint64(1), cmd.ID)
			return domain.Tag{ID: cmd.ID, UserID: cmd.UserID, CategoryID: 2, Name: *cmd.Name, CreatedAt: now, UpdatedAt: now}, nil
		},
		getByIDFn: func(_ context.Context, tagID uint64, userID uint64) (domain.Tag, error) {
			return domain.Tag{ID: tagID, UserID: userID, CategoryID: 2, Name: "ByID", CreatedAt: now, UpdatedAt: now}, nil
		},
		getByNameFn: func(_ context.Context, tagName string, userID uint64) (domain.Tag, error) {
			return domain.Tag{ID: 3, UserID: userID, CategoryID: 2, Name: tagName, CreatedAt: now, UpdatedAt: now}, nil
		},
		getByCategoryFn: func(_ context.Context, categoryID uint64, userID uint64) ([]domain.Tag, error) {
			return []domain.Tag{{ID: 4, UserID: userID, CategoryID: categoryID, Name: "ByCategory", CreatedAt: now, UpdatedAt: now}}, nil
		},
		getAllFn: func(_ context.Context, userID uint64) ([]domain.Tag, error) {
			return []domain.Tag{{ID: 5, UserID: userID, CategoryID: 2, Name: "All", CreatedAt: now, UpdatedAt: now}}, nil
		},
		softDeleteFn: func(_ context.Context, tagID, userID uint64) error {
			require.Equal(t, uint64(1), tagID)
			require.Equal(t, uint64(9), userID)
			return nil
		},
	}, tagLoggerStub{})

	created, err := ctrl.Create(t.Context(), gmodel.CreateTagInput{Name: "Focus", CategoryID: "2", Description: &desc, Icon: &icon}, 9)
	require.NoError(t, err)
	require.Equal(t, "1", created.ID)

	updated, err := ctrl.Update(t.Context(), gmodel.UpdateTagInput{ID: "1", Name: &newName, CategoryID: &newCategoryID}, 9)
	require.NoError(t, err)
	require.Equal(t, "Updated", updated.Name)

	byID, err := ctrl.GetByID(t.Context(), 7, 9)
	require.NoError(t, err)
	require.Equal(t, "7", byID.ID)

	byName, err := ctrl.GetByName(t.Context(), "Focus", 9)
	require.NoError(t, err)
	require.Equal(t, "Focus", byName.Name)

	byCategory, err := ctrl.GetByCategoryID(t.Context(), 2, 9)
	require.NoError(t, err)
	require.Len(t, byCategory, 1)

	all, err := ctrl.GetAll(t.Context(), 9)
	require.NoError(t, err)
	require.Len(t, all, 1)

	err = ctrl.SoftDelete(t.Context(), 1, 9)
	require.NoError(t, err)
}

func TestTagController_GuardsAndErrors(t *testing.T) {
	ctrl := controller.NewController(&tagServiceStub{
		createFn: func(context.Context, taginput.CreateTagCommand) (domain.Tag, error) {
			return domain.Tag{}, errors.New("create failed")
		},
		updateFn: func(context.Context, taginput.UpdateTagCommand) (domain.Tag, error) {
			return domain.Tag{}, errors.New("update failed")
		},
		getByIDFn: func(context.Context, uint64, uint64) (domain.Tag, error) {
			return domain.Tag{}, errors.New("get failed")
		},
		getByNameFn: func(context.Context, string, uint64) (domain.Tag, error) {
			return domain.Tag{}, errors.New("get failed")
		},
		getByCategoryFn: func(context.Context, uint64, uint64) ([]domain.Tag, error) {
			return nil, errors.New("list failed")
		},
		getAllFn: func(context.Context, uint64) ([]domain.Tag, error) {
			return nil, errors.New("list failed")
		},
		softDeleteFn: func(context.Context, uint64, uint64) error {
			return errors.New("delete failed")
		},
	}, tagLoggerStub{})

	_, err := ctrl.Create(t.Context(), gmodel.CreateTagInput{Name: "x", CategoryID: "bad"}, 1)
	require.ErrorIs(t, err, controller.ErrInvalidCategoryID)

	_, err = ctrl.Create(t.Context(), gmodel.CreateTagInput{Name: "x", CategoryID: "1"}, 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)

	_, err = ctrl.Create(t.Context(), gmodel.CreateTagInput{Name: "x", CategoryID: "1"}, 1)
	require.EqualError(t, err, "create failed")

	_, err = ctrl.Update(t.Context(), gmodel.UpdateTagInput{ID: "bad"}, 1)
	require.ErrorIs(t, err, controller.ErrInvalidTagID)

	badCategory := "bad"
	_, err = ctrl.Update(t.Context(), gmodel.UpdateTagInput{ID: "1", CategoryID: &badCategory}, 1)
	require.ErrorIs(t, err, controller.ErrInvalidTagID)

	_, err = ctrl.Update(t.Context(), gmodel.UpdateTagInput{ID: "1"}, 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)

	_, err = ctrl.Update(t.Context(), gmodel.UpdateTagInput{ID: "1"}, 1)
	require.EqualError(t, err, "update failed")

	_, err = ctrl.GetByID(t.Context(), 0, 1)
	require.ErrorIs(t, err, controller.ErrTagNotFound)

	_, err = ctrl.GetByID(t.Context(), 1, 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)

	_, err = ctrl.GetByID(t.Context(), 1, 1)
	require.EqualError(t, err, "get failed")

	_, err = ctrl.GetByName(t.Context(), "", 1)
	require.ErrorIs(t, err, controller.ErrTagNotFound)

	_, err = ctrl.GetByName(t.Context(), "x", 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)

	_, err = ctrl.GetByName(t.Context(), "x", 1)
	require.EqualError(t, err, "get failed")

	_, err = ctrl.GetByCategoryID(t.Context(), 0, 1)
	require.ErrorIs(t, err, controller.ErrCategoryNotFound)

	_, err = ctrl.GetByCategoryID(t.Context(), 1, 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)

	_, err = ctrl.GetByCategoryID(t.Context(), 1, 1)
	require.EqualError(t, err, "list failed")

	_, err = ctrl.GetAll(t.Context(), 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)

	_, err = ctrl.GetAll(t.Context(), 1)
	require.EqualError(t, err, "list failed")

	err = ctrl.SoftDelete(t.Context(), 1, 1)
	require.EqualError(t, err, "delete failed")
}

var _ taginput.TagService = (*tagServiceStub)(nil)
