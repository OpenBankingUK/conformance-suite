package server

import (
	"net/http"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/discovery"
	"github.com/labstack/echo"
)

type discoveryHandlers struct {
	webJourney Journey
}

func newDiscoveryHandlers(webJourney Journey) discoveryHandlers {
	return discoveryHandlers{webJourney}
}

func (d discoveryHandlers) setDiscoveryModelHandler(c echo.Context) error {
	discoveryModel := &discovery.Model{}
	if err := c.Bind(discoveryModel); err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	failures, err := d.webJourney.SetDiscoveryModel(discoveryModel)
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}

	if !failures.Empty() {
		return c.JSON(http.StatusBadRequest, validationFailuresResponse{failures})
	}

	return c.JSON(http.StatusCreated, discoveryModel)
}

type validationFailuresResponse struct {
	Error discovery.ValidationFailures `json:"error"`
}
