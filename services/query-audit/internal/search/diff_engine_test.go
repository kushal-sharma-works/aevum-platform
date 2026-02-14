package search

import (
	"testing"

	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/domain"
)

func TestDiffEngineDiffTracksAddedAndRemoved(t *testing.T) {
	engine := &DiffEngine{}

	t1 := map[string]*domain.IndexedDecision{
		"d1": {DecisionID: "d1"},
		"d2": {DecisionID: "d2"},
	}
	t2 := map[string]*domain.IndexedDecision{
		"d2": {DecisionID: "d2"},
		"d3": {DecisionID: "d3"},
	}

	result := engine.diff(t1, t2)

	if len(result.Removed) != 1 || result.Removed[0] != "d1" {
		t.Fatalf("unexpected removed decisions: %+v", result.Removed)
	}
	if len(result.Added) != 1 || result.Added[0] != "d3" {
		t.Fatalf("unexpected added decisions: %+v", result.Added)
	}
	if result.Summary != "Added: 1, Removed: 1" {
		t.Fatalf("unexpected summary: %s", result.Summary)
	}
}

func TestDiffEngineDiffNoChanges(t *testing.T) {
	engine := &DiffEngine{}

	t1 := map[string]*domain.IndexedDecision{"d1": {DecisionID: "d1"}}
	t2 := map[string]*domain.IndexedDecision{"d1": {DecisionID: "d1"}}

	result := engine.diff(t1, t2)

	if len(result.Added) != 0 || len(result.Removed) != 0 {
		t.Fatalf("expected no changes, got added=%v removed=%v", result.Added, result.Removed)
	}
}
