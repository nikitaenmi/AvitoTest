package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nikitaenmi/AvitoTest/internal/domain"
)

type Handlers struct {
	service domain.Service
}

func NewHandlers(svc domain.Service) *Handlers {
	return &Handlers{service: svc}
}

func (h *Handlers) HealthCheck(c echo.Context) error {
	ctx := c.Request().Context()
	if err := h.service.HealthCheck(ctx); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}
