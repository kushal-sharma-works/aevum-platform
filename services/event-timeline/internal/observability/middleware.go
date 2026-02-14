package observability

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func GinOTelMiddleware(service string) gin.HandlerFunc {
	return otelgin.Middleware(service)
}

func EchoOTelMiddleware(service string) echo.MiddlewareFunc {
	return otelecho.Middleware(service)
}
