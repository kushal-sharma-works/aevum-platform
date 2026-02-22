package search

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/domain"
)

func TestEngineSearchSuccessAndError(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		if strings.Contains(r.URL.Path, "_search") {
			if r.URL.Path == "/aevum-decisions/_search" {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"took":1,"hits":{"total":{"value":1},"hits":[{"_source":{"id":"1"}}]}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	})

	server := httptest.NewServer(h)
	defer server.Close()

	es, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{server.URL}})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	engine := NewEngine(es, slog.New(slog.NewTextHandler(io.Discard, nil)))

	result, err := engine.Search(context.Background(), "payment", "all", "", 0, 10)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if _, ok := result["hits"]; !ok {
		t.Fatal("expected hits in search result")
	}

	if _, err := engine.Search(context.Background(), "payment", "decisions", "", 0, 10); err == nil {
		t.Fatal("expected domain error for 500 search")
	}
}

func TestDiffEngineQueryDecisions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		if strings.Contains(r.URL.Path, "_search") {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"hits":{"hits":[{"_source":{"decision_id":"d1","rule_id":"r1","rule_version":"2","status":"approved","output":{"x":1},"input":{"y":2}}}]}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	es, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{server.URL}})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	engine := NewDiffEngine(es, slog.New(slog.NewTextHandler(io.Discard, nil)))
	decisions, err := engine.queryDecisions(context.Background(), time.Now(), "r1", "2", "s1")
	if err != nil {
		t.Fatalf("expected query success, got %v", err)
	}
	if len(decisions) != 1 {
		t.Fatalf("expected one decision, got %d", len(decisions))
	}
	if decisions["d1"].Status != "approved" {
		t.Fatal("expected mapped decision status")
	}
}

func TestDiffEngineDiffDetectsChangedFields(t *testing.T) {
	engine := &DiffEngine{}
	before := map[string]*domain.IndexedDecision{"d1": {DecisionID: "d1", Status: "rejected", RuleVersion: "1"}}
	after := map[string]*domain.IndexedDecision{"d1": {DecisionID: "d1", Status: "approved", RuleVersion: "2"}}

	res := engine.diff(before, after)
	if len(res.Changed) != 2 {
		t.Fatalf("expected 2 changed fields, got %d", len(res.Changed))
	}
}
