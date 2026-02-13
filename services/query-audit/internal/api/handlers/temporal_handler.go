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

	fromTime, _ := time.Parse(time.RFC3339, from)
	toTime, _ := time.Parse(time.RFC3339, to)

	pageInt, _ := strconv.Atoi(page)
	if pageInt < 1 {
		pageInt = 1
	}
	pageSizeInt, _ := strconv.Atoi(pageSize)
	if pageSizeInt < 1 {
		pageSizeInt = 50
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

	t1, _ := time.Parse(time.RFC3339, t1Str)
	t2, _ := time.Parse(time.RFC3339, t2Str)

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

	trail, err := ah.builder.Build(c, decisionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "audit build failed"})
		return
	}

	c.JSON(http.StatusOK, trail)
}
