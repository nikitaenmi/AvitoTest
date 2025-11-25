package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nikitaenmi/AvitoTest/internal/domain"
)

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (h *Handlers) handleError(c echo.Context, err error) error {
	var domainErr *domain.DomainError
	if ok := errors.As(err, &domainErr); ok {
		return h.handleDomainError(c, domainErr)
	}

	if echoErr, ok := err.(*echo.HTTPError); ok {
		return echoErr
	}

	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		},
	})
}

func (h *Handlers) handleDomainError(c echo.Context, domainErr *domain.DomainError) error {
	var statusCode int

	switch domainErr.Type {
	case domain.ErrorTypeTeamExists, domain.ErrorTypeValidation:
		statusCode = http.StatusBadRequest
	case domain.ErrorTypePRExists, domain.ErrorTypePRMerged, domain.ErrorTypeNotAssigned, domain.ErrorTypeNoCandidate:
		statusCode = http.StatusConflict
	case domain.ErrorTypeNotFound:
		statusCode = http.StatusNotFound
	default:
		statusCode = http.StatusBadRequest
	}

	return c.JSON(statusCode, ErrorResponse{
		Error: struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}{
			Code:    string(domainErr.Type),
			Message: domainErr.Message,
		},
	})
}
