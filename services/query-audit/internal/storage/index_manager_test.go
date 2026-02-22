package storage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestIndexManagerLifecycle(t *testing.T) {
	var mu sync.Mutex
	exists := map[string]bool{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		path := strings.TrimPrefix(r.URL.Path, "/")

		if strings.HasSuffix(path, "_stats") {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"_all":{"primaries":{}}}`))
			return
		}

		index := strings.Split(path, "/")[0]
		mu.Lock()
		defer mu.Unlock()

		switch r.Method {
		case http.MethodHead:
			if exists[index] {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		case http.MethodPut:
			exists[index] = true
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"acknowledged":true}`))
		case http.MethodDelete:
			delete(exists, index)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"acknowledged":true}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	ec, err := NewElasticsearchClient([]string{server.URL})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	im := NewIndexManager(ec.GetClient())

	if err := im.CreateIndexes(context.Background()); err != nil {
		t.Fatalf("expected create indexes success, got %v", err)
	}

	ok, err := im.IndexExists(context.Background(), "aevum-events")
	if err != nil || !ok {
		t.Fatalf("expected existing index, ok=%v err=%v", ok, err)
	}

	if _, err := im.GetIndexStats(context.Background(), "aevum-events"); err != nil {
		t.Fatalf("expected stats success, got %v", err)
	}

	if err := im.DeleteIndex(context.Background(), "aevum-events"); err != nil {
		t.Fatalf("expected delete success, got %v", err)
	}
}

func TestIndexManagerStatsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		if strings.HasSuffix(r.URL.Path, "_stats") {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ec, err := NewElasticsearchClient([]string{server.URL})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	im := NewIndexManager(ec.GetClient())

	if _, err := im.GetIndexStats(context.Background(), "aevum-events"); err == nil {
		t.Fatal("expected stats error")
	}
}
