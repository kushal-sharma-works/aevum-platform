package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/search"
)

// TemporalHandler handles temporal query requests
type TemporalHandler struct {
	query *search.TemporalQuery
}

// NewTemporalHandler creates a new temporal handler
func NewTemporalHandler(q *search.TemporalQuery) *TemporalHandler {
	return &TemporalHandler{query: q}
}

// Handle executes a temporal query
func (th *TemporalHandler) Handle(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	streamID := c.Query("stream_id")
	queryType := c.DefaultQuery("type", "all")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("size", "50")

	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to are required"})
		return
	}
	if queryType != "all" && queryType != "events" && queryType != "decisions" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be one of: all, events, decisions"})
		return
	}

	fromTime, err := time.Parse(time.RFC3339, from)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from must be RFC3339"})
		return
	}
	toTime, err := time.Parse(time.RFC3339, to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to must be RFC3339"})
		return
	}
	if toTime.Before(fromTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "to must be greater than or equal to from"})
		return
	}

	pageInt, _ := strconv.Atoi(page)
	if pageInt < 1 {
		pageInt = 1
	}
	pageSizeInt, _ := strconv.Atoi(pageSize)
	if pageSizeInt < 1 {
		pageSizeInt = 50
	}
	if pageSizeInt > 200 {
		pageSizeInt = 200
	}

	tq := &domain.TemporalQuery{
		StreamID: streamID,
		From:     fromTime,
		To:       toTime,
		Type:     queryType,
		Page:     pageInt,
		Size:     pageSizeInt,
	}

	results, err := th.query.Execute(c, tq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// CorrelationHandler handles correlation query requests
type CorrelationHandler struct {
	query *search.CorrelationQuery
}

// NewCorrelationHandler creates a new correlation handler
func NewCorrelationHandler(q *search.CorrelationQuery) *CorrelationHandler {
	return &CorrelationHandler{query: q}
}

// Handle executes a correlation query
func (ch *CorrelationHandler) Handle(c *gin.Context) {
	eventID := c.Query("event_id")
	decisionID := c.Query("decision_id")
	ruleID := c.Query("rule_id")
	ruleVersionStr := c.Query("rule_version")
	streamID := c.Query("stream_id")
	eventType := c.Query("event_type")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("size", "50")

	pageInt, _ := strconv.Atoi(page)
	if pageInt < 1 {
		pageInt = 1
	}
	pageSizeInt, _ := strconv.Atoi(pageSize)
	if pageSizeInt < 1 {
		pageSizeInt = 50
	}
	if pageSizeInt > 200 {
		pageSizeInt = 200
	}
	if eventID == "" && decisionID == "" && ruleID == "" && streamID == "" && eventType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one correlation filter is required"})
		return
	}

	cq := &domain.CorrelationQuery{
		EventID:     eventID,
		DecisionID:  decisionID,
		RuleID:      ruleID,
		RuleVersion: ruleVersionStr,
		StreamID:    streamID,
		EventType:   eventType,
		Page:        pageInt,
		Size:        pageSizeInt,
	}

	results, err := ch.query.Execute(c, cq)
	if err != nil {
		if domErr, ok := err.(*domain.DomainError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": domErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// DiffHandler handles decision diff requests
type DiffHandler struct {
	engine *search.DiffEngine
}

// NewDiffHandler creates a new diff handler
func NewDiffHandler(e *search.DiffEngine) *DiffHandler {
	return &DiffHandler{engine: e}
}

// Handle executes a diff query
func (dh *DiffHandler) Handle(c *gin.Context) {
	t1Str := c.Query("t1")
	t2Str := c.Query("t2")
	streamID := c.Query("stream_id")
	ruleID := c.Query("rule_id")
	ruleVersion := c.Query("rule_version")
	if ruleVersion == "" {
		ruleVersion = c.Query("v2")
	}

	if t1Str == "" || t2Str == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "t1 and t2 are required"})
		return
	}
	t1, err := time.Parse(time.RFC3339, t1Str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "t1 must be RFC3339"})
		return
	}
	t2, err := time.Parse(time.RFC3339, t2Str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "t2 must be RFC3339"})
		return
	}
	if t2.Before(t1) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "t2 must be greater than or equal to t1"})
		return
	}

	dq := &domain.DiffQuery{
		T1:          t1,
		T2:          t2,
		RuleID:      ruleID,
		RuleVersion: ruleVersion,
		StreamID:    streamID,
	}

	result, err := dh.engine.Compare(c, dq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "diff failed"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// AuditHandler handles audit trail requests
type AuditHandler struct {
	builder *search.AuditBuilder
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler(b *search.AuditBuilder) *AuditHandler {
	return &AuditHandler{builder: b}
}

// Handle builds and returns an audit trail
func (ah *AuditHandler) Handle(c *gin.Context) {
	decisionID := c.Param("decisionId")
	if decisionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "decisionId is required"})
		return
	}

	trail, err := ah.builder.Build(c, decisionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "audit build failed"})
		return
	}

	c.JSON(http.StatusOK, trail)
}
