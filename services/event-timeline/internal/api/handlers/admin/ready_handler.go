package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReadyHandler struct{}

func NewReadyHandler() *ReadyHandler {
	return &ReadyHandler{}
}

func (h *ReadyHandler) GetReady(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ready"})
}
