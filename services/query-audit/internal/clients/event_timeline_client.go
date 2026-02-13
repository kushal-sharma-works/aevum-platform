package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// EventTimelineClient communicates with Event Timeline Service
type EventTimelineClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewEventTimelineClient creates a new client
func NewEventTimelineClient(baseURL string) *EventTimelineClient {
	return &EventTimelineClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetEvents fetches events with cursor-based pagination
func (c *EventTimelineClient) GetEvents(ctx context.Context, cursor string, limit int) ([]map[string]interface{}, string, error) {
	url := fmt.Sprintf("%s/api/v1/streams/all/events?limit=%d", c.baseURL, limit)
	if cursor != "" {
		url += fmt.Sprintf("&cursor=%s", cursor)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to fetch events: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, "", err
	}

	events := []map[string]interface{}{}
	if data, ok := result["data"].([]interface{}); ok {
		for _, e := range data {
			if event, ok := e.(map[string]interface{}); ok {
				events = append(events, event)
			}
		}
	}

	newCursor := cursor
	if meta, ok := result["meta"].(map[string]interface{}); ok {
		if nc, ok := meta["next_cursor"].(string); ok {
			newCursor = nc
		}
	}

	return events, newCursor, nil
}

// GetEvent fetches a single event by ID
func (c *EventTimelineClient) GetEvent(ctx context.Context, eventID string) (*map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/events/%s", c.baseURL, eventID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch event: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		return &data, nil
	}

	return nil, fmt.Errorf("no data in response")
}

// DecisionEngineClient communicates with Decision Engine Service
type DecisionEngineClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewDecisionEngineClient creates a new client
func NewDecisionEngineClient(baseURL string) *DecisionEngineClient {
	return &DecisionEngineClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetDecisions fetches decisions in a time range
func (c *DecisionEngineClient) GetDecisions(ctx context.Context, from, to time.Time) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/decisions?from=%s&to=%s", c.baseURL, from.Format(time.RFC3339), to.Format(time.RFC3339))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch decisions: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	decisions := []map[string]interface{}{}
	if data, ok := result["data"].([]interface{}); ok {
		for _, d := range data {
			if decision, ok := d.(map[string]interface{}); ok {
				decisions = append(decisions, decision)
			}
		}
	}

	return decisions, nil
}

// GetDecision fetches a single decision by ID
func (c *DecisionEngineClient) GetDecision(ctx context.Context, decisionID string) (*map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/decisions/%s", c.baseURL, decisionID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch decision: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		return &data, nil
	}

	return nil, fmt.Errorf("no data in response")
}

// GetRules fetches all rules
func (c *DecisionEngineClient) GetRules(ctx context.Context) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/rules", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch rules: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		return data, nil
	}

	return make(map[string]interface{}), nil
}

// GetRule fetches a single rule by ID
func (c *DecisionEngineClient) GetRule(ctx context.Context, ruleID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/rules/%s", c.baseURL, ruleID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch rule: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		return data, nil
	}

	return make(map[string]interface{}), nil
}
