package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
)

type stubEventRepository struct {
	events             []domain.Event
	listPendingErr     error
	markPublishedCalls []markPublishedCall
	rescheduleCalls    []rescheduleCall
	markPublishedErr   error
	rescheduleErr      error
}

type markPublishedCall struct {
	eventID     string
	publishedAt time.Time
}

type rescheduleCall struct {
	eventID         string
	nextAvailableAt time.Time
	lastError       string
}

func (r *stubEventRepository) Save(context.Context, domain.Event) error {
	return nil
}

func (r *stubEventRepository) ListPending(context.Context, int) ([]domain.Event, error) {
	if r.listPendingErr != nil {
		return nil, r.listPendingErr
	}
	return r.events, nil
}

func (r *stubEventRepository) MarkPublished(_ context.Context, eventID string, publishedAt time.Time) error {
	if r.markPublishedErr != nil {
		return r.markPublishedErr
	}
	r.markPublishedCalls = append(r.markPublishedCalls, markPublishedCall{eventID: eventID, publishedAt: publishedAt})
	return nil
}

func (r *stubEventRepository) Reschedule(_ context.Context, eventID string, nextAvailableAt time.Time, lastError string) error {
	if r.rescheduleErr != nil {
		return r.rescheduleErr
	}
	r.rescheduleCalls = append(r.rescheduleCalls, rescheduleCall{
		eventID:         eventID,
		nextAvailableAt: nextAvailableAt,
		lastError:       lastError,
	})
	return nil
}

type stubEventPublisher struct {
	publishErrByEventID map[string]error
	publishedEventIDs   []string
}

func (p *stubEventPublisher) Publish(_ context.Context, event domain.Event) error {
	if err := p.publishErrByEventID[event.EventID]; err != nil {
		return err
	}
	p.publishedEventIDs = append(p.publishedEventIDs, event.EventID)
	return nil
}

type noopLogger struct{}

func (noopLogger) Infof(string, ...any)                      {}
func (noopLogger) Errorf(string, ...any)                     {}
func (noopLogger) Debugf(string, ...any)                     {}
func (noopLogger) Warnf(string, ...any)                      {}
func (noopLogger) Infow(string, ...any)                      {}
func (noopLogger) Errorw(string, ...any)                     {}
func (noopLogger) Debugw(string, ...any)                     {}
func (noopLogger) Warnw(string, ...any)                      {}
func (noopLogger) InfowCtx(context.Context, string, ...any)  {}
func (noopLogger) ErrorwCtx(context.Context, string, ...any) {}
func (noopLogger) WarnwCtx(context.Context, string, ...any)  {}
func (noopLogger) DebugwCtx(context.Context, string, ...any) {}

func TestPublishPendingMarksPublishedOnSuccess(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.March, 13, 12, 0, 0, 0, time.UTC)
	repo := &stubEventRepository{
		events: []domain.Event{{
			EventID:       "evt-1",
			AggregateType: "record",
			AggregateID:   "42",
			EventType:     "record.created",
		}},
	}
	publisher := &stubEventPublisher{}
	service := &PublisherService{
		repository: repo,
		publisher:  publisher,
		logger:     noopLogger{},
		now: func() time.Time {
			return now
		},
		backoff: DefaultPublishBackoff,
	}

	if err := service.PublishPending(context.Background(), 10); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(publisher.publishedEventIDs) != 1 || publisher.publishedEventIDs[0] != "evt-1" {
		t.Fatalf("expected publisher to emit evt-1, got %#v", publisher.publishedEventIDs)
	}

	if len(repo.markPublishedCalls) != 1 {
		t.Fatalf("expected one mark-published call, got %d", len(repo.markPublishedCalls))
	}
	if repo.markPublishedCalls[0].eventID != "evt-1" {
		t.Fatalf("expected mark-published for evt-1, got %#v", repo.markPublishedCalls[0])
	}
	if !repo.markPublishedCalls[0].publishedAt.Equal(now) {
		t.Fatalf("expected publishedAt %v, got %v", now, repo.markPublishedCalls[0].publishedAt)
	}
	if len(repo.rescheduleCalls) != 0 {
		t.Fatalf("expected no reschedule calls, got %d", len(repo.rescheduleCalls))
	}
}

func TestPublishPendingReschedulesFailedEvents(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, time.March, 13, 12, 30, 0, 0, time.UTC)
	repo := &stubEventRepository{
		events: []domain.Event{{
			EventID:       "evt-2",
			AggregateType: "record",
			AggregateID:   "99",
			EventType:     "record.created",
		}},
	}
	publisher := &stubEventPublisher{
		publishErrByEventID: map[string]error{
			"evt-2": errors.New("kafka unavailable"),
		},
	}
	service := &PublisherService{
		repository: repo,
		publisher:  publisher,
		logger:     noopLogger{},
		now: func() time.Time {
			return now
		},
		backoff: 5 * time.Second,
	}

	err := service.PublishPending(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error")
	}

	if len(repo.markPublishedCalls) != 0 {
		t.Fatalf("expected no mark-published calls, got %d", len(repo.markPublishedCalls))
	}
	if len(repo.rescheduleCalls) != 1 {
		t.Fatalf("expected one reschedule call, got %d", len(repo.rescheduleCalls))
	}
	if repo.rescheduleCalls[0].eventID != "evt-2" {
		t.Fatalf("expected reschedule for evt-2, got %#v", repo.rescheduleCalls[0])
	}
	if !repo.rescheduleCalls[0].nextAvailableAt.Equal(now.Add(5 * time.Second)) {
		t.Fatalf("expected next available at %v, got %v", now.Add(5*time.Second), repo.rescheduleCalls[0].nextAvailableAt)
	}
	if repo.rescheduleCalls[0].lastError != "kafka unavailable" {
		t.Fatalf("expected original publish error, got %q", repo.rescheduleCalls[0].lastError)
	}
}

func TestPublishPendingReturnsListPendingError(t *testing.T) {
	t.Parallel()

	repo := &stubEventRepository{listPendingErr: errors.New("db down")}
	service := &PublisherService{
		repository: repo,
		publisher:  &stubEventPublisher{},
		logger:     noopLogger{},
		now:        time.Now,
		backoff:    DefaultPublishBackoff,
	}

	err := service.PublishPending(context.Background(), 10)
	if err == nil || err.Error() != "db down" {
		t.Fatalf("expected db down error, got %v", err)
	}
}
