package server

import (
	"net/http"

	"github.com/labstack/echo"
)

type runHandlers struct {
	webJourney Journey
}

// POST /api/run/start
func (h *runHandlers) runStartPostHandler(c echo.Context) error {
	report, err := h.webJourney.RunTests()
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	return c.JSON(http.StatusCreated, report)
}
