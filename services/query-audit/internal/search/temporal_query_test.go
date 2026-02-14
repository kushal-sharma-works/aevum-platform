package search

import (
	"testing"
)

func TestParseResultsHasMoreAndNextPage(t *testing.T) {
	esResult := map[string]interface{}{
		"took": float64(12),
		"hits": map[string]interface{}{
			"total": map[string]interface{}{"value": float64(3)},
			"hits": []interface{}{
				map[string]interface{}{"_source": map[string]interface{}{"event_id": "e1"}},
				map[string]interface{}{"_source": map[string]interface{}{"event_id": "e2"}},
			},
		},
	}

	res := parseResults(esResult, 0, 2)
	if res.Total != 3 {
		t.Fatalf("expected total 3, got %d", res.Total)
	}
	if len(res.Hits) != 2 {
		t.Fatalf("expected 2 hits, got %d", len(res.Hits))
	}
	if !res.HasMore {
		t.Fatal("expected HasMore=true")
	}
	if res.NextPage != 1 {
		t.Fatalf("expected next page 1, got %d", res.NextPage)
	}
}

func TestParseResultsNoMorePages(t *testing.T) {
	esResult := map[string]interface{}{
		"took": float64(8),
		"hits": map[string]interface{}{
			"total": map[string]interface{}{"value": float64(2)},
			"hits": []interface{}{
				map[string]interface{}{"_source": map[string]interface{}{"event_id": "e1"}},
				map[string]interface{}{"_source": map[string]interface{}{"event_id": "e2"}},
			},
		},
	}

	res := parseResults(esResult, 0, 2)
	if res.HasMore {
		t.Fatal("expected HasMore=false")
	}
}
