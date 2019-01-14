package server

import (
	"net/http"

	"github.com/labstack/echo"
)

type reportingEndpoints struct {
	webJourney Journey
}

func newReportingEndpoints(webJourney Journey) reportingEndpoints {
	return reportingEndpoints{webJourney}
}

func (d reportingEndpoints) handler(c echo.Context) error {
	result, err := d.webJourney.RunTests()
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}
	return c.JSON(http.StatusOK, result)
}
