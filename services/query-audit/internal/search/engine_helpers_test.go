package search

import (
	"encoding/json"
	"testing"
)

func TestSearchIndexes(t *testing.T) {
	tests := map[string][]string{
		"events":    {"aevum-events"},
		"decisions": {"aevum-decisions"},
		"all":       {"aevum-events", "aevum-decisions"},
		"":          {"aevum-events", "aevum-decisions"},
	}

	for input, expected := range tests {
		indexes := searchIndexes(input)
		if len(indexes) != len(expected) {
			t.Fatalf("expected %d indexes, got %d", len(expected), len(indexes))
		}
		for i := range expected {
			if indexes[i] != expected[i] {
				t.Fatalf("expected index %s, got %s", expected[i], indexes[i])
			}
		}
	}
}

func TestBuildSearchQueryContainsPagingAndFilter(t *testing.T) {
	body := buildSearchQuery("payment", "all", "stream-1", 20, 10)

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		t.Fatalf("failed to parse query body: %v", err)
	}

	if int(parsed["from"].(float64)) != 20 {
		t.Fatalf("expected from=20")
	}
	if int(parsed["size"].(float64)) != 10 {
		t.Fatalf("expected size=10")
	}
}

func TestTemporalIndexesAndParseResultsMalformed(t *testing.T) {
	if len(temporalIndexes("events")) != 1 {
		t.Fatal("expected events index")
	}
	if len(temporalIndexes("all")) != 2 {
		t.Fatal("expected two indexes for all")
	}

	res := parseResults(map[string]interface{}{}, 1, 50)
	if res == nil {
		t.Fatal("expected non-nil results")
	}
	if len(res.Hits) != 0 {
		t.Fatal("expected empty hits for malformed payload")
	}
}
