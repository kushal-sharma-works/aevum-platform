package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/replay"
)

type adminEventStore struct {
	latest int64
	err    error
	events []domain.Event
}

func (s *adminEventStore) PutEvent(context.Context, domain.Event) error         { return nil }
func (s *adminEventStore) PutEventsBatch(context.Context, []domain.Event) error { return nil }
func (s *adminEventStore) GetByEventID(context.Context, string) (domain.Event, error) {
	return domain.Event{}, domain.ErrNotFound
}
func (s *adminEventStore) FindByIdempotencyKey(context.Context, string, string) (domain.Event, error) {
	return domain.Event{}, domain.ErrNotFound
}
func (s *adminEventStore) GetLatestSequence(context.Context, string) (int64, error) {
	return s.latest, s.err
}
func (s *adminEventStore) QueryByStream(_ context.Context, streamID string, from int64, _ string, _ int32) ([]domain.Event, int64, bool, error) {
	filtered := make([]domain.Event, 0)
	for _, event := range s.events {
		if event.StreamID == streamID && event.SequenceNumber >= from {
			filtered = append(filtered, event)
		}
	}
	return filtered, from + int64(len(filtered)), false, nil
}

type adminStreamStore struct {
	streams []domain.Stream
	err     error
}

func (s *adminStreamStore) ListStreams(context.Context, int32) ([]domain.Stream, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.streams, nil
}

func TestHealthHandlerStatuses(t *testing.T) {
	e := echo.New()
	h := NewHealthHandler(&adminEventStore{latest: 1})
	req := httptest.NewRequest(http.MethodGet, "/admin/health", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := h.GetHealth(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "ok")

	hDegraded := NewHealthHandler(&adminEventStore{err: errors.New("db down")})
	rec2 := httptest.NewRecorder()
	ctx2 := e.NewContext(httptest.NewRequest(http.MethodGet, "/admin/health", nil), rec2)
	err = hDegraded.GetHealth(ctx2)
	require.NoError(t, err)
	require.Contains(t, rec2.Body.String(), "degraded")
}

func TestReadyHandler(t *testing.T) {
	e := echo.New()
	h := NewReadyHandler()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/admin/ready", nil), rec)

	err := h.GetReady(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "ready")
}

func TestStreamsHandler(t *testing.T) {
	e := echo.New()
	h := NewStreamsHandler(&adminStreamStore{streams: []domain.Stream{{StreamID: "s1", LatestSequence: 2}}})
	rec := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/admin/streams", nil), rec)

	err := h.ListStreams(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "s1")

	hErr := NewStreamsHandler(&adminStreamStore{err: errors.New("boom")})
	recErr := httptest.NewRecorder()
	ctxErr := e.NewContext(httptest.NewRequest(http.MethodGet, "/admin/streams", nil), recErr)
	err = hErr.ListStreams(ctxErr)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, recErr.Code)
}

func TestReplayHandler(t *testing.T) {
	event, err := domain.NewEvent(domain.NewEventInput{
		EventID:        "evt-admin",
		StreamID:       "stream-admin",
		SequenceNumber: 1,
		EventType:      "created",
		Payload:        json.RawMessage(`{"k":"v"}`),
		OccurredAt:     time.Now().UTC(),
		IngestedAt:     time.Now().UTC(),
		SchemaVersion:  1,
	})
	require.NoError(t, err)

	store := &adminEventStore{events: []domain.Event{event}}
	engine := replay.NewEngine(store, observability.NewMetrics())
	h := NewReplayHandler(engine)

	e := echo.New()
	body := `{"stream_id":"stream-admin","from":"2026-02-14T00:00:00Z","to":"2026-02-15T00:00:00Z","event_types":["created"],"page_size":50}`
	req := httptest.NewRequest(http.MethodPost, "/admin/replay", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err = h.TriggerReplay(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "completed")
}

func TestMetricsHandler(t *testing.T) {
	e := echo.New()
	h := NewMetricsHandler(observability.NewMetrics())
	rec := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/admin/metrics", nil), rec)

	err := h.GetMetrics(ctx)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "aevum_active_replays")
}
