package api

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/clients"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/search"
)

func TestSetupRouter_BasicEndpointsAndMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter(nil, nil, nil, nil, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if w.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected request id header")
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/admin/metrics", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestSetupRouter_RegistersSearchRoutesWhenDependenciesProvided(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	searchEngine := search.NewEngine(nil, logger)
	temporal := search.NewTemporalQuery(nil, logger)
	correlation := search.NewCorrelationQuery(nil, logger)
	diff := search.NewDiffEngine(nil, logger)
	audit := search.NewAuditBuilder(nil, clients.NewEventTimelineClient("http://example"), clients.NewDecisionEngineClient("http://example"), logger)

	router := SetupRouter(searchEngine, temporal, correlation, diff, audit)

	routeSet := map[string]bool{}
	for _, route := range router.Routes() {
		routeSet[route.Method+" "+route.Path] = true
	}

	required := []string{
		http.MethodGet + " /api/v1/search",
		http.MethodPost + " /api/v1/search",
		http.MethodGet + " /api/v1/timeline",
		http.MethodPost + " /api/v1/timeline",
		http.MethodGet + " /api/v1/correlate",
		http.MethodPost + " /api/v1/correlate",
		http.MethodGet + " /api/v1/diff",
		http.MethodPost + " /api/v1/diff",
		http.MethodGet + " /api/v1/audit/:decisionId",
	}

	for _, key := range required {
		if !routeSet[key] {
			t.Fatalf("expected route %s to be registered", key)
		}
	}
}
