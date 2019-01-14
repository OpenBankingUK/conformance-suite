package server

import (
	"net/http"

	"github.com/labstack/echo"
)

type testCaseHandlers struct {
	webJourney Journey
}

func newTestCaseHandlers(webJourney Journey) testCaseHandlers {
	return testCaseHandlers{webJourney}
}

func (d testCaseHandlers) testCasesHandler(c echo.Context) error {
	testCases, err := d.webJourney.TestCases()
	if err != nil {
		return c.JSON(http.StatusBadRequest, NewErrorResponse(err))
	}
	return c.JSON(http.StatusOK, testCases)
}
