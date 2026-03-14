package controller_test

import (
	"context"
	"errors"
	"testing"
	"time"

	gmodel "github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/category/adapter/primary/graphql/controller"
	"github.com/lechitz/AionApi/internal/category/core/domain"
	catinput "github.com/lechitz/AionApi/internal/category/core/ports/input"
	"github.com/stretchr/testify/require"
)

type categoryServiceStub struct {
	createFn     func(context.Context, catinput.CreateCategoryCommand) (domain.Category, error)
	updateFn     func(context.Context, catinput.UpdateCategoryCommand) (domain.Category, error)
	getByIDFn    func(context.Context, uint64, uint64) (domain.Category, error)
	getByNameFn  func(context.Context, string, uint64) (domain.Category, error)
	listAllFn    func(context.Context, uint64) ([]domain.Category, error)
	softDeleteFn func(context.Context, uint64, uint64) error
}

func (s *categoryServiceStub) Create(ctx context.Context, cmd catinput.CreateCategoryCommand) (domain.Category, error) {
	if s.createFn == nil {
		panic("unexpected Create call")
	}
	return s.createFn(ctx, cmd)
}

func (s *categoryServiceStub) Update(ctx context.Context, cmd catinput.UpdateCategoryCommand) (domain.Category, error) {
	if s.updateFn == nil {
		panic("unexpected Update call")
	}
	return s.updateFn(ctx, cmd)
}

func (s *categoryServiceStub) GetByID(ctx context.Context, categoryID, userID uint64) (domain.Category, error) {
	if s.getByIDFn == nil {
		panic("unexpected GetByID call")
	}
	return s.getByIDFn(ctx, categoryID, userID)
}

func (s *categoryServiceStub) GetByName(ctx context.Context, categoryName string, userID uint64) (domain.Category, error) {
	if s.getByNameFn == nil {
		panic("unexpected GetByName call")
	}
	return s.getByNameFn(ctx, categoryName, userID)
}

func (s *categoryServiceStub) ListAll(ctx context.Context, userID uint64) ([]domain.Category, error) {
	if s.listAllFn == nil {
		panic("unexpected ListAll call")
	}
	return s.listAllFn(ctx, userID)
}

func (s *categoryServiceStub) SoftDelete(ctx context.Context, categoryID, userID uint64) error {
	if s.softDeleteFn == nil {
		panic("unexpected SoftDelete call")
	}
	return s.softDeleteFn(ctx, categoryID, userID)
}

type categoryLoggerStub struct{}

func (categoryLoggerStub) Infof(string, ...any)                      {}
func (categoryLoggerStub) Errorf(string, ...any)                     {}
func (categoryLoggerStub) Debugf(string, ...any)                     {}
func (categoryLoggerStub) Warnf(string, ...any)                      {}
func (categoryLoggerStub) Infow(string, ...any)                      {}
func (categoryLoggerStub) Errorw(string, ...any)                     {}
func (categoryLoggerStub) Debugw(string, ...any)                     {}
func (categoryLoggerStub) Warnw(string, ...any)                      {}
func (categoryLoggerStub) InfowCtx(context.Context, string, ...any)  {}
func (categoryLoggerStub) ErrorwCtx(context.Context, string, ...any) {}
func (categoryLoggerStub) WarnwCtx(context.Context, string, ...any)  {}
func (categoryLoggerStub) DebugwCtx(context.Context, string, ...any) {}

func TestCategoryController_Basics(t *testing.T) {
	now := time.Date(2026, 2, 14, 12, 0, 0, 0, time.UTC)
	desc := "desc"
	color := "#fff"
	icon := "star"

	ctrl := controller.NewController(&categoryServiceStub{
		createFn: func(_ context.Context, cmd catinput.CreateCategoryCommand) (domain.Category, error) {
			require.Equal(t, uint64(10), cmd.UserID)
			require.Equal(t, "Work", cmd.Name)
			return domain.Category{ID: 1, UserID: cmd.UserID, Name: cmd.Name, Description: desc, Color: color, Icon: icon, CreatedAt: now, UpdatedAt: now}, nil
		},
		updateFn: func(_ context.Context, cmd catinput.UpdateCategoryCommand) (domain.Category, error) {
			require.Equal(t, uint64(1), cmd.ID)
			return domain.Category{ID: cmd.ID, UserID: cmd.UserID, Name: "Updated", CreatedAt: now, UpdatedAt: now}, nil
		},
		getByIDFn: func(_ context.Context, categoryID, userID uint64) (domain.Category, error) {
			return domain.Category{ID: categoryID, UserID: userID, Name: "ByID", CreatedAt: now, UpdatedAt: now}, nil
		},
		getByNameFn: func(_ context.Context, categoryName string, userID uint64) (domain.Category, error) {
			return domain.Category{ID: 2, UserID: userID, Name: categoryName, CreatedAt: now, UpdatedAt: now}, nil
		},
		listAllFn: func(_ context.Context, userID uint64) ([]domain.Category, error) {
			return []domain.Category{{ID: 3, UserID: userID, Name: "All", CreatedAt: now, UpdatedAt: now}}, nil
		},
		softDeleteFn: func(_ context.Context, categoryID, userID uint64) error {
			require.Equal(t, uint64(9), categoryID)
			require.Equal(t, uint64(10), userID)
			return nil
		},
	}, categoryLoggerStub{})

	created, err := ctrl.Create(t.Context(), gmodel.CreateCategoryInput{Name: "Work", Description: &desc, ColorHex: &color, Icon: &icon}, 10)
	require.NoError(t, err)
	require.Equal(t, "1", created.ID)
	require.NotNil(t, created.Description)

	updated, err := ctrl.Update(t.Context(), gmodel.UpdateCategoryInput{ID: "1", Name: &[]string{"Updated"}[0]}, 10)
	require.NoError(t, err)
	require.Equal(t, "Updated", updated.Name)

	byID, err := ctrl.GetByID(t.Context(), 5, 10)
	require.NoError(t, err)
	require.Equal(t, "5", byID.ID)

	byName, err := ctrl.GetByName(t.Context(), "Gym", 10)
	require.NoError(t, err)
	require.Equal(t, "Gym", byName.Name)

	all, err := ctrl.ListAll(t.Context(), 10)
	require.NoError(t, err)
	require.Len(t, all, 1)

	err = ctrl.SoftDelete(t.Context(), 9, 10)
	require.NoError(t, err)
}

func TestCategoryController_GuardsAndErrors(t *testing.T) {
	ctrl := controller.NewController(&categoryServiceStub{
		createFn: func(context.Context, catinput.CreateCategoryCommand) (domain.Category, error) {
			return domain.Category{}, errors.New("create failed")
		},
		updateFn: func(context.Context, catinput.UpdateCategoryCommand) (domain.Category, error) {
			return domain.Category{}, errors.New("update failed")
		},
		getByIDFn: func(context.Context, uint64, uint64) (domain.Category, error) {
			return domain.Category{}, errors.New("get failed")
		},
		getByNameFn: func(context.Context, string, uint64) (domain.Category, error) {
			return domain.Category{}, errors.New("get failed")
		},
		listAllFn: func(context.Context, uint64) ([]domain.Category, error) {
			return nil, errors.New("list failed")
		},
		softDeleteFn: func(context.Context, uint64, uint64) error {
			return errors.New("delete failed")
		},
	}, categoryLoggerStub{})

	_, err := ctrl.Create(t.Context(), gmodel.CreateCategoryInput{Name: "x"}, 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)

	_, err = ctrl.Create(t.Context(), gmodel.CreateCategoryInput{Name: "x"}, 1)
	require.EqualError(t, err, "create failed")

	_, err = ctrl.Update(t.Context(), gmodel.UpdateCategoryInput{ID: "bad"}, 1)
	require.ErrorIs(t, err, controller.ErrInvalidCategoryID)

	_, err = ctrl.Update(t.Context(), gmodel.UpdateCategoryInput{ID: "1"}, 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)

	_, err = ctrl.Update(t.Context(), gmodel.UpdateCategoryInput{ID: "1"}, 1)
	require.EqualError(t, err, "update failed")

	_, err = ctrl.GetByID(t.Context(), 0, 1)
	require.ErrorIs(t, err, controller.ErrCategoryNotFound)

	_, err = ctrl.GetByID(t.Context(), 1, 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)

	_, err = ctrl.GetByID(t.Context(), 1, 1)
	require.EqualError(t, err, "get failed")

	_, err = ctrl.GetByName(t.Context(), "", 1)
	require.ErrorIs(t, err, controller.ErrCategoryNotFound)

	_, err = ctrl.GetByName(t.Context(), "x", 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)

	_, err = ctrl.GetByName(t.Context(), "x", 1)
	require.EqualError(t, err, "get failed")

	_, err = ctrl.ListAll(t.Context(), 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)

	_, err = ctrl.ListAll(t.Context(), 1)
	require.EqualError(t, err, "list failed")

	err = ctrl.SoftDelete(t.Context(), 0, 1)
	require.ErrorIs(t, err, controller.ErrCategoryIDNotFound)

	err = ctrl.SoftDelete(t.Context(), 1, 0)
	require.ErrorIs(t, err, controller.ErrUserIDNotFound)

	err = ctrl.SoftDelete(t.Context(), 1, 1)
	require.EqualError(t, err, "delete failed")
}

var _ catinput.CategoryService = (*categoryServiceStub)(nil)
