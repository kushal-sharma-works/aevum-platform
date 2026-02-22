package httputil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestWriteErrorIncludesRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Set("request_id", "req-123")

	WriteError(c, http.StatusBadRequest, "bad_request", "invalid payload")

	require.Equal(t, http.StatusBadRequest, recorder.Code)

	var body ErrorBody
	err := json.Unmarshal(recorder.Body.Bytes(), &body)
	require.NoError(t, err)
	require.Equal(t, "bad_request", body.Error.Code)
	require.Equal(t, "invalid payload", body.Error.Message)
	require.Equal(t, "req-123", body.Error.RequestID)
}

func TestErrorHelpersUseExpectedStatusCodes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name   string
		call   func(*gin.Context)
		status int
	}{
		{name: "bad request", call: func(c *gin.Context) { BadRequest(c, "code", "msg") }, status: http.StatusBadRequest},
		{name: "unauthorized", call: func(c *gin.Context) { Unauthorized(c, "code", "msg") }, status: http.StatusUnauthorized},
		{name: "too many", call: func(c *gin.Context) { TooManyRequests(c, "code", "msg") }, status: http.StatusTooManyRequests},
		{name: "internal", call: func(c *gin.Context) { Internal(c, "code", "msg") }, status: http.StatusInternalServerError},
		{name: "not found", call: func(c *gin.Context) { NotFound(c, "code", "msg") }, status: http.StatusNotFound},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			testCase.call(c)
			require.Equal(t, testCase.status, recorder.Code)
		})
	}
}
