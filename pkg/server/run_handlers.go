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
	result := map[string]interface{}{}
	result["status"] = "executing"

	report, err := h.webJourney.RunTests()
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}
	_ = report
	return c.JSON(http.StatusCreated, result)
}
