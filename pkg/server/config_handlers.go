package server

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
)

type configHandlers struct {
	webJourney Journey
}

type GlobalConfiguration struct {
	SigningPrivate   string `json:"signing_private"`
	SigningPublic    string `json:"signing_public"`
	TransportPrivate string `json:"transport_private"`
	TransportPublic  string `json:"transport_public"`
}

// POST /api/config/global
func (h *configHandlers) configGlobalPostHandler(c echo.Context) error {
	globalConfiguration := new(GlobalConfiguration)
	if err := c.Bind(globalConfiguration); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			NewErrorResponse(errors.Wrap(err, "error with Bind")),
		)
	}

	certificateSigning, err := authentication.NewCertificate(
		globalConfiguration.SigningPublic,
		globalConfiguration.SigningPrivate,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			NewErrorResponse(errors.Wrap(err, "error with signing certificate")),
		)
	}

	certificateTransport, err := authentication.NewCertificate(
		globalConfiguration.TransportPublic,
		globalConfiguration.TransportPrivate,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			NewErrorResponse(errors.Wrap(err, "error with transport certificate")),
		)
	}
	h.webJourney.SetCertificates(certificateSigning, certificateTransport)

	return c.JSON(http.StatusCreated, globalConfiguration)
}
