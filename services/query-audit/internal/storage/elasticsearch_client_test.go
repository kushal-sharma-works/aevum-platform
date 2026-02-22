package storage

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestElasticsearchClientHealthAndGetClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		if r.URL.Path == "/_cluster/health" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"green"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	ec, err := NewElasticsearchClient([]string{server.URL})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if ec.GetClient() == nil {
		t.Fatal("expected underlying client")
	}

	if err := ec.Health(context.Background()); err != nil {
		t.Fatalf("expected healthy cluster, got %v", err)
	}
}

func TestElasticsearchClientHealthReturnsErrorOnNon200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"status":"red"}`))
	}))
	defer server.Close()

	ec, err := NewElasticsearchClient([]string{server.URL})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if err := ec.Health(context.Background()); err == nil {
		t.Fatal("expected health error on non-200")
	}
}
