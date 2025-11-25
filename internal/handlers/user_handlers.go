package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nikitaenmi/AvitoTest/internal/handlers/dto"
)

func (h *Handlers) SetUserActive(c echo.Context) error {
	var req dto.SetUserActiveRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	ctx := c.Request().Context()
	if err := h.service.SetUserActive(ctx, req.ToUserFilter(), req.IsActive); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user activity updated"})
}

func (h *Handlers) GetUserReviewPRs(c echo.Context) error {
	userID := c.QueryParam("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id is required"})
	}

	ctx := c.Request().Context()
	prs, err := h.service.GetUserReviewPRs(ctx, dto.UserFilterFromQuery(userID))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"pull_requests": prs})
}
