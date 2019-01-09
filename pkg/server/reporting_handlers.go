package server

import (
	"net/http"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/web"
	"github.com/labstack/echo"
)

type reportingEndpoints struct {
	webJourney web.Journey
}

func newReportingEndpoints(webJourney web.Journey) reportingEndpoints {
	return reportingEndpoints{webJourney}
}

func (d reportingEndpoints) handler(c echo.Context) error {
	result, err := d.webJourney.RunTests()
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}
	return c.JSON(http.StatusOK, result)
}
