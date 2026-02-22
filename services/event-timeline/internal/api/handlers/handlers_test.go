package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/ingest"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/pkg/clock"
)

type testEventStore struct {
	events        []domain.Event
	byIdem        map[string]domain.Event
	getByEventErr error
	queryErr      error
}

func (s *testEventStore) PutEvent(_ context.Context, event domain.Event) error {
	s.events = append(s.events, event)
	if s.byIdem == nil {
		s.byIdem = map[string]domain.Event{}
	}
	if event.IdempotencyKey != "" {
		s.byIdem[event.StreamID+"#"+event.IdempotencyKey] = event
	}
	return nil
}

func (s *testEventStore) PutEventsBatch(_ context.Context, events []domain.Event) error {
	s.events = append(s.events, events...)
	for _, event := range events {
		if event.IdempotencyKey != "" {
			if s.byIdem == nil {
				s.byIdem = map[string]domain.Event{}
			}
			s.byIdem[event.StreamID+"#"+event.IdempotencyKey] = event
		}
	}
	return nil
}

func (s *testEventStore) GetByEventID(_ context.Context, eventID string) (domain.Event, error) {
	if s.getByEventErr != nil {
		return domain.Event{}, s.getByEventErr
	}
	for _, event := range s.events {
		if event.EventID == eventID {
			return event, nil
		}
	}
	return domain.Event{}, domain.ErrNotFound
}

func (s *testEventStore) FindByIdempotencyKey(_ context.Context, streamID, key string) (domain.Event, error) {
	event, ok := s.byIdem[streamID+"#"+key]
	if !ok {
		return domain.Event{}, domain.ErrNotFound
	}
	return event, nil
}

func (s *testEventStore) GetLatestSequence(_ context.Context, streamID string) (int64, error) {
	latest := int64(0)
	for _, event := range s.events {
		if event.StreamID == streamID && event.SequenceNumber > latest {
			latest = event.SequenceNumber
		}
	}
	return latest, nil
}

func (s *testEventStore) QueryByStream(_ context.Context, streamID string, fromSequence int64, _ string, limit int32) ([]domain.Event, int64, bool, error) {
	if s.queryErr != nil {
		return nil, 0, false, s.queryErr
	}
	filtered := make([]domain.Event, 0)
	for _, event := range s.events {
		if event.StreamID == streamID && event.SequenceNumber >= fromSequence {
			filtered = append(filtered, event)
		}
	}
	if len(filtered) == 0 {
		return filtered, fromSequence, false, nil
	}
	if int32(len(filtered)) > limit {
		filtered = filtered[:limit]
	}
	next := filtered[len(filtered)-1].SequenceNumber + 1
	return filtered, next, false, nil
}

type fixedGenerator struct{}

func (fixedGenerator) New(time.Time) (string, error) { return "01HFIXEDID", nil }

type captureStreamStore struct {
	fromSequence int64
	direction    string
}

type failingPutStore struct{}

func (f *failingPutStore) PutEvent(context.Context, domain.Event) error { return errors.New("put failed") }
func (f *failingPutStore) PutEventsBatch(context.Context, []domain.Event) error { return nil }
func (f *failingPutStore) GetByEventID(context.Context, string) (domain.Event, error) {
	return domain.Event{}, domain.ErrNotFound
}
func (f *failingPutStore) FindByIdempotencyKey(context.Context, string, string) (domain.Event, error) {
	return domain.Event{}, domain.ErrNotFound
}
func (f *failingPutStore) GetLatestSequence(context.Context, string) (int64, error) { return 0, nil }
func (f *failingPutStore) QueryByStream(context.Context, string, int64, string, int32) ([]domain.Event, int64, bool, error) {
	return nil, 0, false, nil
}

func (s *captureStreamStore) PutEvent(context.Context, domain.Event) error         { return nil }
func (s *captureStreamStore) PutEventsBatch(context.Context, []domain.Event) error { return nil }
func (s *captureStreamStore) GetByEventID(context.Context, string) (domain.Event, error) {
	return domain.Event{}, domain.ErrNotFound
}
func (s *captureStreamStore) FindByIdempotencyKey(context.Context, string, string) (domain.Event, error) {
	return domain.Event{}, domain.ErrNotFound
}
func (s *captureStreamStore) GetLatestSequence(context.Context, string) (int64, error) {
	return 5, nil
}
func (s *captureStreamStore) QueryByStream(_ context.Context, _ string, fromSequence int64, direction string, _ int32) ([]domain.Event, int64, bool, error) {
	s.fromSequence = fromSequence
	s.direction = direction
	return []domain.Event{}, fromSequence, false, nil
}

func newIngestService(store *testEventStore) *ingest.Service {
	return ingest.NewService(store, fixedGenerator{}, clock.MockClock{Current: time.Now().UTC()}, observability.NewMetrics())
}

func TestIngestHandlerIngestSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &testEventStore{byIdem: map[string]domain.Event{}}
	handler := NewIngestHandler(newIngestService(store))

	r := gin.New()
	r.POST("/events", handler.Ingest)

	body := `{"stream_id":"stream-1","event_type":"created","payload":{"v":1},"occurred_at":"2026-02-14T10:00:00Z","idempotency_key":"idem-1"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.Contains(t, rec.Body.String(), `"created":true`)
}

func TestIngestHandlerRejectsInvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewIngestHandler(newIngestService(&testEventStore{byIdem: map[string]domain.Event{}}))
	r := gin.New()
	r.POST("/events", handler.Ingest)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestIngestHandlerReturnsInternalOnStoreError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := ingest.NewService(&failingPutStore{}, fixedGenerator{}, clock.MockClock{Current: time.Now().UTC()}, observability.NewMetrics())
	handler := NewIngestHandler(service)
	r := gin.New()
	r.POST("/events", handler.Ingest)

	body := `{"stream_id":"stream-1","event_type":"created","payload":{"v":1},"occurred_at":"2026-02-14T10:00:00Z","idempotency_key":"idem-1"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(rec, req)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestBatchIngestHandlerRejectsInvalidBatchSize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewBatchIngestHandler(newIngestService(&testEventStore{byIdem: map[string]domain.Event{}}))
	r := gin.New()
	r.POST("/events/batch", handler.IngestBatch)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/events/batch", strings.NewReader("[]"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEventHandlerGetByIDNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewEventHandler(&testEventStore{events: nil})
	r := gin.New()
	r.GET("/events/:eventId", handler.GetByID)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/events/unknown", nil))

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestStreamHandlerGetByStreamInvalidLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewStreamHandler(&testEventStore{})
	r := gin.New()
	r.GET("/streams/:streamId/events", handler.GetByStream)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/streams/s1/events?limit=abc", nil))

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestStreamHandlerGetByStreamInvalidDirection(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewStreamHandler(&testEventStore{})
	r := gin.New()
	r.GET("/streams/:streamId/events", handler.GetByStream)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/streams/s1/events?direction=sideways", nil))

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestStreamHandlerGetByStreamInvalidCursor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewStreamHandler(&testEventStore{})
	r := gin.New()
	r.GET("/streams/:streamId/events", handler.GetByStream)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/streams/s1/events?cursor=invalid", nil))

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEventHandlerGetByIDInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewEventHandler(&testEventStore{getByEventErr: errors.New("db unavailable")})
	r := gin.New()
	r.GET("/events/:eventId", handler.GetByID)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/events/evt-1", nil))

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestStreamHandlerGetByStreamSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	event, err := domain.NewEvent(domain.NewEventInput{
		EventID:        "evt-1",
		StreamID:       "s1",
		SequenceNumber: 1,
		EventType:      "created",
		Payload:        json.RawMessage(`{"x":1}`),
		OccurredAt:     time.Now().UTC(),
		IngestedAt:     time.Now().UTC(),
		SchemaVersion:  1,
	})
	require.NoError(t, err)
	store := &testEventStore{events: []domain.Event{event}}

	handler := NewStreamHandler(store)
	r := gin.New()
	r.GET("/streams/:streamId/events", handler.GetByStream)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/streams/s1/events?limit=10", nil))

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "events")
}

func TestStreamHandlerBackwardDefaultsToLatestSequence(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := &captureStreamStore{}
	handler := NewStreamHandler(store)
	r := gin.New()
	r.GET("/streams/:streamId/events", handler.GetByStream)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/streams/s1/events?direction=backward", nil))

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, int64(5), store.fromSequence)
	require.Equal(t, domain.DirectionBackward, store.direction)
}
