package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nikitaenmi/AvitoTest/internal/handlers/dto"
)

func (h *Handlers) CreatePR(c echo.Context) error {
	var req dto.CreatePRRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	ctx := c.Request().Context()
	pr, err := h.service.CreatePR(ctx, req.ToDomain())
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"pr": pr})
}

func (h *Handlers) MergePR(c echo.Context) error {
	var req dto.MergePRRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	ctx := c.Request().Context()
	if err := h.service.MergePR(ctx, req.ToPRFilter()); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "pull request merged"})
}

func (h *Handlers) ReassignReviewer(c echo.Context) error {
	var req dto.ReassignReviewerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	ctx := c.Request().Context()

	newReviewerID, err := h.service.ReassignReviewer(ctx, req.ToPRFilter(), req.OldUserID)
	if err != nil {
		return h.handleError(c, err)
	}

	updatedPR, err := h.service.GetPR(ctx, req.ToPRFilter())
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"pr":          updatedPR,
		"replaced_by": newReviewerID,
	})
}
