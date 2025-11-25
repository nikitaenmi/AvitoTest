package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nikitaenmi/AvitoTest/internal/handlers/dto"
)

func (h *Handlers) CreateTeam(c echo.Context) error {
	var req dto.CreateTeamRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
	}

	ctx := c.Request().Context()
	team, err := h.service.CreateTeam(ctx, req.ToDomain())
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"team": team})
}

func (h *Handlers) GetTeam(c echo.Context) error {
	teamName := c.QueryParam("team_name")
	if teamName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "team_name is required"})
	}

	ctx := c.Request().Context()
	team, err := h.service.GetTeam(ctx, dto.TeamFilterFromQuery(teamName))
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"team": team})
}
