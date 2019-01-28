package server

import (
	"bitbucket.org/openbankingteam/conformance-suite/proxy"
	"context"
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"bitbucket.org/openbankingteam/conformance-suite/appconfig"
	"bitbucket.org/openbankingteam/conformance-suite/pkg/authentication"

	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

type configHandlers struct {
	server     *Server
	webJourney Journey
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

// createProxy - kick off proxy by:
// loading the spec,
// creating a new proxy configured with
//    - bind address
//    - swagger specification location
//    - target host (aspsp resource server)
//    - verbosity
// configure an default logreport
func createProxy(appConfig *appconfig.AppConfig) (*http.Server, error) {
	logrus.Info("Server:createProxy -> Proxy")

	appConfig.PrintAppConfig()
	doc, err := loads.Spec(appConfig.Spec)
	if err != nil {
		logrus.Errorln("Server:createProxy -> loads.Spec err=", err)
		return nil, err
	}

	proxy, err := proxy.New(
		doc.Spec(),
		&proxy.LogReporter{},
		proxy.WithTarget(appConfig.TargetHost),
		proxy.WithVerbose(appConfig.Verbose),
		proxy.WithAppConfig(appConfig),
	)
	if err != nil {
		logrus.Errorln("Server:createProxy -> proxy.New err=", err)
		return nil, err
	}

	// start serving the proxy - and don't return unless there is a problem/exit
	// also sleep for a bit until it starts...
	server, serveErr := serveProxy(proxy, appConfig.Bind)
	time.Sleep(200 * time.Millisecond)

	// block until serveErr has an error value or the specified timeout has elapsed.
	// we might need to bump up this timeout to something a bit larger.
	timeout := time.After(1 * time.Second)
	select {
	case err := <-serveErr: // Error from listen&serve - exit
		return nil, err
	case <-timeout:
	}

	logrus.WithFields(logrus.Fields{
		"bind":   appConfig.Bind,
		"target": proxy.Target(),
	}).Info("Server:createProxy -> Proxy is listening")

	// Report PendingOperations - part of shutdown tidyup
	logrus.Debugln("Pending Operations:")
	for i, op := range proxy.PendingOperations() {
		logrus.Debugf("%03d) id=%s", i+1, op.ID)
	}

	return server, nil
}

// Run the proxy at the address specified by "bind"
// Requests get sent to the target server identified by proxy.Target()
// configure some channels to handle shutdown/interrupts
//
// Return channel so that caller can block waiting
// to see if we managed to start the server or not.
func serveProxy(proxy *proxy.Proxy, bind string) (*http.Server, chan error) {
	server := &http.Server{
		Addr:    bind,
		Handler: proxy.Router(),
	}

	// Run server s.ListenAndServe on a goroutine
	serveErr := make(chan error)
	go func() {
		err := server.ListenAndServe()
		serveErr <- err
	}()

	return server, serveErr
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
	if err := h.server.proxy.Shutdown(context.TODO()); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			NewErrorResponse(err),
		)
	}

	h.server.proxy = nil
	h.server.logger.Debugf("Server:configDeleteHandler -> status=down proxy=%+v", h.server.proxy)

	return c.NoContent(http.StatusOK)
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
