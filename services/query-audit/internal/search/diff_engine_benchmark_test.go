package search

import (
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/domain"
)

func BenchmarkDiffEngineDiff(b *testing.B) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	engine := &DiffEngine{logger: logger}

	t1 := make(map[string]*domain.IndexedDecision, 10000)
	t2 := make(map[string]*domain.IndexedDecision, 10000)

	for i := 0; i < 10000; i++ {
		id := fmt.Sprintf("decision-%d", i)
		t1[id] = &domain.IndexedDecision{DecisionID: id}
	}
	for i := 5000; i < 15000; i++ {
		id := fmt.Sprintf("decision-%d", i)
		t2[id] = &domain.IndexedDecision{DecisionID: id}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := engine.diff(t1, t2)
		if len(res.Added) != 5000 || len(res.Removed) != 5000 {
			b.Fatalf("unexpected diff result: added=%d removed=%d", len(res.Added), len(res.Removed))
		}
	}
}
