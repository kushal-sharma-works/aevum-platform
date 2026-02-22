package clients

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// EventTimelineClient communicates with Event Timeline Service
type EventTimelineClient struct {
	baseURL    string
	jwtSecret  string
	defaultStream string
	httpClient *http.Client
}

// NewEventTimelineClient creates a new client
func NewEventTimelineClient(baseURL string) *EventTimelineClient {
	return &EventTimelineClient{
		baseURL: baseURL,
		jwtSecret: os.Getenv("EVENT_TIMELINE_JWT_SECRET"),
		defaultStream: func() string {
			stream := os.Getenv("EVENT_TIMELINE_DEFAULT_STREAM")
			if strings.TrimSpace(stream) == "" {
				return "default"
			}
			return stream
		}(),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *EventTimelineClient) addAuth(req *http.Request) {
	if strings.TrimSpace(c.jwtSecret) == "" {
		return
	}
	now := time.Now().Unix()
	exp := now + 3600
	headerJSON := `{"alg":"HS256","typ":"JWT"}`
	payloadJSON := fmt.Sprintf(`{"iss":"query-audit","sub":"sync-worker","iat":%d,"exp":%d}`, now, exp)
	header := base64.RawURLEncoding.EncodeToString([]byte(headerJSON))
	payload := base64.RawURLEncoding.EncodeToString([]byte(payloadJSON))
	unsigned := header + "." + payload
	h := hmac.New(sha256.New, []byte(c.jwtSecret))
	_, _ = h.Write([]byte(unsigned))
	sig := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	req.Header.Set("Authorization", "Bearer "+unsigned+"."+sig)
}

// GetEvents fetches events with cursor-based pagination
func (c *EventTimelineClient) GetEvents(ctx context.Context, cursor string, limit int) ([]map[string]interface{}, string, error) {
	endpoint := fmt.Sprintf("%s/api/v1/streams/%s/events?limit=%d", c.baseURL, c.defaultStream, limit)
	if cursor != "" {
		endpoint += fmt.Sprintf("&cursor=%s", url.QueryEscape(cursor))
	}

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, "", err
	}
	c.addAuth(req)

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
	if data, ok := result["events"].([]interface{}); ok {
		for _, e := range data {
			if event, ok := e.(map[string]interface{}); ok {
				events = append(events, event)
			}
		}
	}

	newCursor := cursor
	if nc, ok := result["next_cursor"].(string); ok {
		newCursor = nc
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
	c.addAuth(req)

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

	return &result, nil
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
	decisions := []map[string]interface{}{}

	rulesURL := fmt.Sprintf("%s/api/v1/rules", c.baseURL)
	rulesReq, err := http.NewRequestWithContext(ctx, "GET", rulesURL, nil)
	if err != nil {
		return nil, err
	}

	rulesResp, err := c.httpClient.Do(rulesReq)
	if err != nil {
		return nil, err
	}
	defer rulesResp.Body.Close()

	if rulesResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch decisions: status %d", rulesResp.StatusCode)
	}

	rulesBody, err := io.ReadAll(rulesResp.Body)
	if err != nil {
		return nil, err
	}

	var rules []map[string]interface{}
	if err := json.Unmarshal(rulesBody, &rules); err != nil {
		return nil, err
	}

	for _, rule := range rules {
		ruleID, _ := rule["id"].(string)
		if strings.TrimSpace(ruleID) == "" {
			continue
		}

		decisionsURL := fmt.Sprintf("%s/api/v1/decisions/rule/%s", c.baseURL, ruleID)
		decisionsReq, err := http.NewRequestWithContext(ctx, "GET", decisionsURL, nil)
		if err != nil {
			continue
		}

		decisionsResp, err := c.httpClient.Do(decisionsReq)
		if err != nil {
			continue
		}

		if decisionsResp.StatusCode != http.StatusOK {
			_ = decisionsResp.Body.Close()
			continue
		}

		decisionBody, err := io.ReadAll(decisionsResp.Body)
		_ = decisionsResp.Body.Close()
		if err != nil {
			continue
		}

		var decisionList []map[string]interface{}
		if err := json.Unmarshal(decisionBody, &decisionList); err != nil {
			continue
		}

		for _, decision := range decisionList {
			timestamp, _ := decision["evaluated_at"].(string)
			if strings.TrimSpace(timestamp) != "" {
				if ts, err := time.Parse(time.RFC3339, timestamp); err == nil {
					if ts.Before(from) || ts.After(to) {
						continue
					}
				}
			}
			decisions = append(decisions, decision)
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

	return &result, nil
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

	var result any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if rules, ok := result.([]interface{}); ok {
		return map[string]interface{}{"rules": rules}, nil
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

	return result, nil
}
