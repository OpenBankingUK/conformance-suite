package server

import (
	"net/http"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/web"
	"github.com/labstack/echo"
)

type testCaseHandlers struct {
	webJourney web.Journey
}

func newTestCaseHandlers(webJourney web.Journey) testCaseHandlers {
	return testCaseHandlers{webJourney}
}

func (d testCaseHandlers) testCasesHandler(c echo.Context) error {
	testCases, err := d.webJourney.TestCases()
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}
	return c.JSON(http.StatusOK, testCases)
}
