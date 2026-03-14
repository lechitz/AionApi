package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/record/core/domain"
)

type projectionRepositoryStub struct {
	item  domain.RecordProjection
	items []domain.RecordProjection
}

func (p projectionRepositoryStub) GetProjectedByID(context.Context, uint64, uint64) (domain.RecordProjection, error) {
	return p.item, nil
}

func (p projectionRepositoryStub) ListProjectedLatest(context.Context, uint64, int) ([]domain.RecordProjection, error) {
	return p.items, nil
}

func (p projectionRepositoryStub) ListProjectedPage(context.Context, uint64, int, *string, *int64) ([]domain.RecordProjection, error) {
	return p.items, nil
}

func TestService_GetProjectedByID(t *testing.T) {
	t.Parallel()

	svc := &Service{
		RecordProjectionRepository: projectionRepositoryStub{
			item: domain.RecordProjection{RecordID: 5177, LastEventType: "record.created"},
		},
	}

	got, err := svc.GetProjectedByID(context.Background(), 5177, 7)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.RecordID != 5177 {
		t.Fatalf("expected record id 5177, got %d", got.RecordID)
	}
}

func TestService_ListProjectedLatestRequiresRepository(t *testing.T) {
	t.Parallel()

	svc := &Service{}
	_, err := svc.ListProjectedLatest(context.Background(), 7, 5)
	if !errors.Is(err, ErrProjectionRepositoryUnavailable) {
		t.Fatalf("expected ErrProjectionRepositoryUnavailable, got %v", err)
	}
}

func TestService_ListProjectedPageDefaultsLimit(t *testing.T) {
	t.Parallel()

	svc := &Service{
		RecordProjectionRepository: projectionRepositoryStub{
			items: []domain.RecordProjection{{RecordID: 1}},
		},
	}

	got, err := svc.ListProjectedPage(context.Background(), 7, 0, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected one projection, got %d", len(got))
	}
}
