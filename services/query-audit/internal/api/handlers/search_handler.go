package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/search"
)

// SearchHandler handles general search requests
type SearchHandler struct {
	engine *search.Engine
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(engine *search.Engine) *SearchHandler {
	return &SearchHandler{engine: engine}
}

// Handle executes a search query
func (sh *SearchHandler) Handle(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	queryType := c.DefaultQuery("type", "all")
	streamID := c.Query("stream_id")
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

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query is required", "code": string(domain.ErrInvalidQuery)})
		return
	}
	if queryType != "all" && queryType != "events" && queryType != "decisions" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be one of: all, events, decisions", "code": string(domain.ErrInvalidQuery)})
		return
	}

	sq := &domain.SearchQuery{
		Query:    query,
		Type:     queryType,
		StreamID: streamID,
		Page:     pageInt,
		Size:     pageSizeInt,
	}

	from := (pageInt - 1) * pageSizeInt
	results, err := sh.engine.Search(c, sq.Query, sq.Type, sq.StreamID, from, sq.Size)
	if err != nil {
		if domErr, ok := err.(*domain.DomainError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": domErr.Message, "code": string(domErr.Code)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, results)
}
