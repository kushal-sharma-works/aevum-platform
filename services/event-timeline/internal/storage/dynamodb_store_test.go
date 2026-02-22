package storage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
)

func testDynamoClient(t *testing.T, handler http.Handler) (*dynamodb.Client, func()) {
	t.Helper()
	server := httptest.NewServer(handler)

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("eu-central-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{URL: server.URL, SigningRegion: region, HostnameImmutable: true}, nil
		})),
	)
	require.NoError(t, err)

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.RetryMaxAttempts = 1
	})

	return client, server.Close
}

func dynamoHandler(status int, response string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-amz-json-1.0")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(response))
	})
}

func sampleEvent(t *testing.T) domain.Event {
	t.Helper()
	event, err := domain.NewEvent(domain.NewEventInput{
		EventID:        "evt-1",
		StreamID:       "stream-1",
		SequenceNumber: 1,
		EventType:      "created",
		Payload:        json.RawMessage(`{"ok":true}`),
		OccurredAt:     time.Date(2026, 2, 14, 12, 0, 0, 0, time.UTC),
		IngestedAt:     time.Date(2026, 2, 14, 12, 0, 0, 0, time.UTC),
		SchemaVersion:  1,
	})
	require.NoError(t, err)
	return event
}

func TestDynamoDBEventStorePutEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, cleanup := testDynamoClient(t, dynamoHandler(http.StatusOK, `{}`))
		defer cleanup()

		store := NewDynamoDBEventStore(client, "events")
		err := store.PutEvent(context.Background(), sampleEvent(t))
		require.NoError(t, err)
	})

	t.Run("sequence conflict", func(t *testing.T) {
		body := `{"__type":"com.amazonaws.dynamodb.v20120810#TransactionCanceledException","CancellationReasons":[{"Code":"ConditionalCheckFailed"},{"Code":"None"}],"message":"Transaction cancelled"}`
		client, cleanup := testDynamoClient(t, dynamoHandler(http.StatusBadRequest, body))
		defer cleanup()

		store := NewDynamoDBEventStore(client, "events")
		err := store.PutEvent(context.Background(), sampleEvent(t))
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrSequenceConflict)
	})
}

func TestDynamoDBEventStoreQueries(t *testing.T) {
	itemResponse := `{"Items":[{"PK":{"S":"evt-1"},"SK":{"S":"EVENT#stream-1#00000000000000000001"},"GSI1PK":{"S":"stream-1"},"GSI1SK":{"N":"1"},"EventType":{"S":"created"},"Payload":{"B":"e30="},"OccurredAt":{"S":"2026-02-14T12:00:00Z"},"IngestedAt":{"S":"2026-02-14T12:00:00Z"},"SchemaVersion":{"N":"1"}}],"LastEvaluatedKey":{"PK":{"S":"evt-1"}}}`
	queryHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-amz-json-1.0")
		_ = r.ParseForm()
		target := r.Header.Get("X-Amz-Target")
		w.WriteHeader(http.StatusOK)
		switch target {
		case "DynamoDB_20120810.Query":
			_, _ = w.Write([]byte(itemResponse))
		default:
			_, _ = w.Write([]byte(`{}`))
		}
	})

	client, cleanup := testDynamoClient(t, queryHandler)
	defer cleanup()
	store := NewDynamoDBEventStore(client, "events")

	event, err := store.GetByEventID(context.Background(), "evt-1")
	require.NoError(t, err)
	require.Equal(t, "evt-1", event.EventID)

	idem, err := store.FindByIdempotencyKey(context.Background(), "stream-1", "idem-1")
	require.NoError(t, err)
	require.Equal(t, "evt-1", idem.EventID)

	seq, err := store.GetLatestSequence(context.Background(), "stream-1")
	require.NoError(t, err)
	require.Equal(t, int64(1), seq)

	events, nextSeq, hasMore, err := store.QueryByStream(context.Background(), "stream-1", 1, domain.DirectionForward, 10)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, int64(2), nextSeq)
	require.True(t, hasMore)
}

func TestDynamoDBEventStoreBatchAndEdgeCases(t *testing.T) {
	t.Run("batch empty", func(t *testing.T) {
		client, cleanup := testDynamoClient(t, dynamoHandler(http.StatusOK, `{}`))
		defer cleanup()

		store := NewDynamoDBEventStore(client, "events")
		err := store.PutEventsBatch(context.Background(), nil)
		require.NoError(t, err)
	})

	t.Run("batch success and failure", func(t *testing.T) {
		calls := 0
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-amz-json-1.0")
			calls++
			if calls == 1 {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"__type":"InternalServerError","message":"boom"}`))
		})

		client, cleanup := testDynamoClient(t, handler)
		defer cleanup()
		store := NewDynamoDBEventStore(client, "events")

		err := store.PutEventsBatch(context.Background(), []domain.Event{sampleEvent(t)})
		require.NoError(t, err)

		err = store.PutEventsBatch(context.Background(), []domain.Event{sampleEvent(t)})
		require.Error(t, err)
	})

	t.Run("query not found and backward paging", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/x-amz-json-1.0")
			target := r.Header.Get("X-Amz-Target")
			w.WriteHeader(http.StatusOK)
			switch target {
			case "DynamoDB_20120810.Query":
				_, _ = w.Write([]byte(`{"Items":[{"PK":{"S":"evt-5"},"SK":{"S":"EVENT#stream-1#00000000000000000005"},"GSI1PK":{"S":"stream-1"},"GSI1SK":{"N":"5"}}]}`))
			default:
				_, _ = w.Write([]byte(`{"Items":[]}`))
			}
		})

		client, cleanup := testDynamoClient(t, handler)
		defer cleanup()
		store := NewDynamoDBEventStore(client, "events")

		events, nextSeq, hasMore, err := store.QueryByStream(context.Background(), "stream-1", 5, domain.DirectionBackward, 10)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, int64(4), nextSeq)
		require.False(t, hasMore)
	})
}

func TestDynamoDBStreamStoreListStreams(t *testing.T) {
	resp := `{"Items":[{"GSI1PK":{"S":"stream-1"},"GSI1SK":{"N":"1"}},{"GSI1PK":{"S":"stream-1"},"GSI1SK":{"N":"3"}},{"GSI1PK":{"S":"stream-2"},"GSI1SK":{"N":"2"}}]}`
	client, cleanup := testDynamoClient(t, dynamoHandler(http.StatusOK, resp))
	defer cleanup()

	store := NewDynamoDBStreamStore(client, "events")
	streams, err := store.ListStreams(context.Background(), 0)
	require.NoError(t, err)
	require.Len(t, streams, 2)
}

func TestDynamoDBEventStoreNotFoundAndErrors(t *testing.T) {
	notFoundHandler := dynamoHandler(http.StatusOK, `{"Items":[]}`)
	client, cleanup := testDynamoClient(t, notFoundHandler)
	defer cleanup()

	store := NewDynamoDBEventStore(client, "events")

	_, err := store.GetByEventID(context.Background(), "missing")
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrNotFound)

	_, err = store.FindByIdempotencyKey(context.Background(), "stream-1", "missing")
	require.Error(t, err)
	require.ErrorIs(t, err, domain.ErrNotFound)

	latest, err := store.GetLatestSequence(context.Background(), "stream-1")
	require.NoError(t, err)
	require.Equal(t, int64(0), latest)

	errorBody := `{"__type":"InternalServerError","message":"boom"}`
	errorClient, errorCleanup := testDynamoClient(t, dynamoHandler(http.StatusInternalServerError, errorBody))
	defer errorCleanup()

	errorStore := NewDynamoDBEventStore(errorClient, "events")
	_, _, _, err = errorStore.QueryByStream(context.Background(), "stream-1", 1, domain.DirectionForward, 10)
	require.Error(t, err)
}

func TestDynamoDBStreamStoreError(t *testing.T) {
	client, cleanup := testDynamoClient(t, dynamoHandler(http.StatusInternalServerError, `{"__type":"InternalServerError","message":"boom"}`))
	defer cleanup()

	store := NewDynamoDBStreamStore(client, "events")
	_, err := store.ListStreams(context.Background(), 10)
	require.Error(t, err)
}
