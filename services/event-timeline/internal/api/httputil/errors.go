package httputil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Error ErrorEnvelope `json:"error"`
}

type ErrorEnvelope struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func WriteError(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorBody{
		Error: ErrorEnvelope{
			Code:      code,
			Message:   message,
			RequestID: c.GetString("request_id"),
		},
	})
}

func BadRequest(c *gin.Context, code, message string) {
	WriteError(c, http.StatusBadRequest, code, message)
}

func Unauthorized(c *gin.Context, code, message string) {
	WriteError(c, http.StatusUnauthorized, code, message)
}

func TooManyRequests(c *gin.Context, code, message string) {
	WriteError(c, http.StatusTooManyRequests, code, message)
}

func Internal(c *gin.Context, code, message string) {
	WriteError(c, http.StatusInternalServerError, code, message)
}

func NotFound(c *gin.Context, code, message string) {
	WriteError(c, http.StatusNotFound, code, message)
}
