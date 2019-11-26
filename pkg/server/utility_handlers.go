package server

import (
	"net/http"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/version"
	"github.com/pkg/errors"

	"github.com/labstack/echo"
)

// utilityEndpoints contains various endpoints that provide utility but don't warrant the
// creation of their own collection of endpoints.
type utilityEndpoints struct {
	version version.Checker
}

func newUtilityEndpoints(version version.Checker) utilityEndpoints {
	return utilityEndpoints{
		version: version,
	}
}

// version is an endpoint that uses the `version` package to determine if an update is available
// for this application
func (u utilityEndpoints) versionCheck(c echo.Context) error {
	vf, err := u.version.VersionFormatter(version.FullVersion)
	if err != nil {
		err = errors.Wrap(err, "format version")
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(err))
	}

	msg, update, err := u.version.UpdateWarningVersion(vf)
	if err != nil {
		err = errors.Wrap(err, "update warning version")
		return c.JSON(http.StatusInternalServerError, NewErrorResponse(err))
	}

	response := VersionResponse{
		Version: u.version.GetHumanVersion(),
		Msg:     msg,
		Update:  update,
	}

	return c.JSON(http.StatusOK, response)
}

// VersionResponse is defined as a response object for /version API calls
type VersionResponse struct {
	Version string `json:"version"`
	Msg     string `json:"message"`
	Update  bool   `json:"update"`
}
