package clients

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEventTimelineClientGetEventsAndGetEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/streams/all/events":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":[{"event_id":"e1"}],"meta":{"next_cursor":"c2"}}`))
		case "/api/v1/events/e1":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"event_id":"e1","stream_id":"s1"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := NewEventTimelineClient(server.URL)

	events, cursor, err := client.GetEvents(context.Background(), "", 50)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(events) != 1 || cursor != "c2" {
		t.Fatalf("unexpected events response: len=%d cursor=%s", len(events), cursor)
	}

	event, err := client.GetEvent(context.Background(), "e1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if (*event)["event_id"] != "e1" {
		t.Fatal("unexpected event payload")
	}
}

func TestDecisionEngineClientEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/decisions":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":[{"decision_id":"d1"}]}`))
		case r.URL.Path == "/api/v1/decisions/d1":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"decision_id":"d1"}}`))
		case r.URL.Path == "/api/v1/rules":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"items":[]}}`))
		case r.URL.Path == "/api/v1/rules/r1":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":{"id":"r1"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := NewDecisionEngineClient(server.URL)

	decisions, err := client.GetDecisions(context.Background(), time.Now().Add(-time.Hour), time.Now())
	if err != nil || len(decisions) != 1 {
		t.Fatalf("unexpected decisions response: err=%v len=%d", err, len(decisions))
	}

	decision, err := client.GetDecision(context.Background(), "d1")
	if err != nil || (*decision)["decision_id"] != "d1" {
		t.Fatalf("unexpected decision response: err=%v value=%v", err, decision)
	}

	rules, err := client.GetRules(context.Background())
	if err != nil || rules == nil {
		t.Fatalf("unexpected rules response: err=%v", err)
	}

	rule, err := client.GetRule(context.Background(), "r1")
	if err != nil || rule["id"] != "r1" {
		t.Fatalf("unexpected rule response: err=%v value=%v", err, rule)
	}
}

func TestClientHandlesNon200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	}))
	defer server.Close()

	eventClient := NewEventTimelineClient(server.URL)
	_, _, err := eventClient.GetEvents(context.Background(), "", 10)
	if err == nil {
		t.Fatal("expected error on non-200 get events")
	}

	decisionClient := NewDecisionEngineClient(server.URL)
	_, err = decisionClient.GetRules(context.Background())
	if err == nil {
		t.Fatal("expected error on non-200 get rules")
	}

	if fmt.Sprint(err) == "" {
		t.Fatal("expected non-empty error")
	}
}
