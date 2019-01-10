package server

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/web"

	"bitbucket.org/openbankingteam/conformance-suite/appconfig"

	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

type configHandlers struct {
	server     *Server
	webJourney web.Journey
}

// POST /api/config
func (h *configHandlers) configPostHandler(c echo.Context) error {
	appConfig := new(appconfig.AppConfig)
	if err := c.Bind(appConfig); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err),
		)
	}
	if err := c.Validate(appConfig); err != nil {
		// translate all error at once
		errs := err.(validator.ValidationErrors)
		errsMap := errs.Translate(nil)

		return c.JSON(http.StatusBadRequest, errsMap)
	}

	h.server.logger.Debugf("Server:configPostHandler -> status=creating proxy")
	proxy, err := createProxy(appConfig)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err),
		)
	}
	h.server.proxy = proxy

	h.server.logger.Debugf("Server:configPostHandler -> status=created proxy=%+v", h.server.proxy)

	return c.JSON(http.StatusOK, appConfig)
}

// DELETE /api/config
func (h *configHandlers) configDeleteHandler(c echo.Context) error {
	if h.server.proxy == nil {
		return c.JSON(
			http.StatusBadRequest,
			NewErrorResponse(fmt.Errorf("proxy has not been configured")),
		)
	}

	h.server.logger.Debugf("Server:configDeleteHandler -> status=destroying down proxy=%+v", h.server.proxy)
	if err := h.server.proxy.Shutdown(nil); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err),
		)
	}

	h.server.proxy = nil
	h.server.logger.Debugf("Server:configDeleteHandler -> status=down proxy=%+v", h.server.proxy)

	return c.NoContent(http.StatusOK)
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
	h.webJourney.SetCertificateSigning(certificateSigning)

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
	h.webJourney.SetCertificateTransport(certificateTransport)

	return c.JSON(http.StatusCreated, globalConfiguration)
}
