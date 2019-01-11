package server

import (
	"net/http"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/web"
	"github.com/labstack/echo"
)

type runHandlers struct {
	webJourney web.Journey
}

// POST /api/run/start
func (h *runHandlers) runStartPostHandler(c echo.Context) error {
	result := map[string]interface{}{}
	result["status"] = "executing"

	// TODO: do something with certificates ...
	h.webJourney.CertificateSigning()
	h.webJourney.CertificateTransport()

	return c.JSON(http.StatusCreated, result)
}
