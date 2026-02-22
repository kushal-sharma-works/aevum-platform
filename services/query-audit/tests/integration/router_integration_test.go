package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/api"
)

func TestRouterHealthEndpoint(t *testing.T) {
	router := api.SetupRouter(nil, nil, nil, nil, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected X-Request-ID header")
	}
}

func TestRouterMetricsEndpoint(t *testing.T) {
	router := api.SetupRouter(nil, nil, nil, nil, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/metrics", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
