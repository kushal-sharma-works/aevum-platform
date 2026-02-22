package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSearchHandlerValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewSearchHandler(nil)
	r := gin.New()
	r.GET("/search", h.Handle)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/search", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing query, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/search?q=test&type=bad", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid type, got %d", w.Code)
	}
}

func TestTemporalHandlerValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewTemporalHandler(nil)
	r := gin.New()
	r.GET("/timeline", h.Handle)

	cases := []string{
		"/timeline",
		"/timeline?from=bad&to=2026-01-01T00:00:00Z",
		"/timeline?from=2026-01-01T00:00:00Z&to=bad",
		"/timeline?from=2026-01-02T00:00:00Z&to=2026-01-01T00:00:00Z",
		"/timeline?from=2026-01-01T00:00:00Z&to=2026-01-02T00:00:00Z&type=bad",
	}

	for _, path := range cases {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for %s, got %d", path, w.Code)
		}
	}
}

func TestCorrelationAndDiffValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	correlation := NewCorrelationHandler(nil)
	diff := NewDiffHandler(nil)

	r := gin.New()
	r.GET("/correlate", correlation.Handle)
	r.GET("/diff", diff.Handle)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/correlate", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty correlation filters, got %d", w.Code)
	}

	diffCases := []string{
		"/diff",
		"/diff?t1=bad&t2=2026-01-01T00:00:00Z",
		"/diff?t1=2026-01-02T00:00:00Z&t2=2026-01-01T00:00:00Z",
	}

	for _, path := range diffCases {
		w = httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for %s, got %d", path, w.Code)
		}
	}
}

func TestAuditHandlerValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewAuditHandler(nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "decisionId", Value: ""}}

	h.Handle(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing decision id, got %d", w.Code)
	}
}
